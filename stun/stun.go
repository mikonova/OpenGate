package stun

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"slices"
	"time"

	"github.com/mikonova/OpenGate/errdef"
	sterr "github.com/mikonova/OpenGate/errdef/stunerrors"
	stattr "github.com/mikonova/OpenGate/stun/stunattr"
	"github.com/mikonova/OpenGate/stun/stunlist"
)

type addrInfo struct {
	IsIPv6  bool
	Port    string
	Address string
}

type attrInfo struct {
	AttrType   int
	AttrValue  []byte
	FullLength uint16
	Padding    uint16
}

var magicCookie []byte = []byte{0x21, 0x12, 0xA4, 0x42}

func StunDial() {
	for {
		if success := dialingLoop(); success {
			break
		}
	}
}

func dialingLoop() (dialSuccess bool) {
	var conn net.Conn
	transactionID := make([]byte, 12)
	inputBuf := make([]byte, 0)
	buf := []byte{
		0x00, 0x01, // binding request
		0x00, 0x00, // message length
	}
	buf = append(buf, magicCookie...)

	_, err := rand.Read(transactionID)
	if err != nil {
		log.Fatalln(err)
	}

	buf = append(buf, transactionID...)

	for _, v := range stunlist.ServerList {
	RETRY:
		conn, err = net.Dial("udp", v)
		if err != nil {
			conn.Close()
			log.Println(errdef.WarnBase + "dialing error on server \"" + v + "\", switching server")
			continue
		}
		conn.SetDeadline(time.Now().Add(time.Second * 3))
		if _, err := conn.Write(buf); err != nil {
			log.Println(errdef.WarnBase + "writing timeout on server \"" + v + "\", switching server")
			conn.Close()
			continue
		}
		if _, err := conn.Read(inputBuf); err != nil {
			log.Println(errdef.WarnBase + "reading timeout on server \"" + v + "\", switching server")
			conn.Close()
			continue
		}

		err, n := decodeResp(inputBuf, transactionID)
		if err != nil && n == sterr.Alternate {
			continue
		} else if err != nil && n == sterr.Retry {
			conn.Close()
			goto RETRY
		} else if err != nil && n == sterr.Other {
			conn.Close()
			continue
		} else if err == nil && n == sterr.Success {
			conn.Close()
			return true
		} else {
			conn.Close()
			continue
		}

	}
	return false

}

/*
err != nil if any sort of error occured. "n" represents what kind of action should be taken.
n = 0 means there's no error, n = 1 means that server should not be switched, n = 2 means alternate server should be tried
n = 3 means any other error, n = 4 means an error is supposed to be fatal (should not happen once configured)
*/
func decodeResp(buf []byte, tID []byte) (err error, n int) {
	argument := attrInfo{}
	header := buf[:20]
	args := buf[20:]
	msgLength := binary.BigEndian.Uint16(header[2:4])
	cookie := header[4:8]
	transactionID := header[8:]

	if n := slices.Compare(transactionID, tID); n != 0 {
		err = errors.New(errdef.ErrBase + "unmatching STUN transactionID")
		log.Println(err.Error())
		return err, sterr.Alternate
	}

	if n := slices.Compare(cookie, magicCookie); n != 0 {
		err = errors.New(errdef.ErrBase + "unmatching STUN magic cookie")
		log.Println(err.Error())
		return err, sterr.Alternate
	}

	for i := 0; i < int(msgLength); {
		args, argument = extractArg(args)
		i += int(argument.FullLength)
		err, n = parseArgument(argument.AttrValue, header, argument.AttrType)
	}
	return err, n

}

func extractArg(argList []byte) (args []byte, argument attrInfo) {
	argType := argList[:2]
	argLen := binary.BigEndian.Uint16(argList[2:4])
	var argumentType int

	if n := slices.Compare(argType, []byte{0x00, 0x20}); n != 0 {
		argumentType = stattr.XORMappedAddr
	} else if n := slices.Compare(argType, []byte{0x00, 0x01}); n != 0 {
		argumentType = stattr.MappedAddr
	} else if n := slices.Compare(argType, []byte{0x00, 0x09}); n != 0 {
		argumentType = stattr.Error
	} else if n := slices.Compare(argType, []byte{0x00, 0x08}); n != 0 {
		argumentType = stattr.MessageIntegrity
	} else {
		argumentType = stattr.Unimportant
	}

	padding := 4 - ((4 + argLen) % 4)

	info := attrInfo{
		AttrType:   argumentType,
		AttrValue:  argList[4:argLen],
		FullLength: 4 + argLen + padding,
		Padding:    padding,
	}
	argList = argList[:argLen+padding]
	return argList, info

}

// TODO: доделать парсинг IP
func parseArgument(argument []byte, header []byte, argtype int) (err error, n int) {
	successResp := []byte{0x01, 0x01}
	errResp := []byte{0x01, 0x11}
	isIPv6 := false

	reqType := header[:2] // 0x01, 0x01 if success response, 0x01, 0x11 if an error response

	if n := slices.Compare(reqType, successResp); n == 0 && argtype == stattr.XORMappedAddr {
		_ = argument[:1]                                                                          // always 0x00
		family := argument[1:2]                                                                   // IPv4 or IPv6
		port := binary.BigEndian.Uint16(argument[2:4]) ^ binary.BigEndian.Uint16(magicCookie[:2]) // MAPPED
		xIP := argument[4:]                                                                       // XOR-MAPPED
		if family[0] == byte(0x01) {
			isIPv6 = true
		}
	} else if slices.Compare(reqType, errResp); n == 0 {
		_ = argument[:2]
		class := binary.BigEndian.Uint16(argument[2:3])
		number := binary.BigEndian.Uint16(argument[3:5])
		reason := string(argument[5:])

		declineStatus := class*100 + number
		switch declineStatus {
		case 300:
			err = errors.New(errdef.ErrBase + "try alternate STUN server: ")
			log.Println(err.Error(), declineStatus, ", reason: ", reason)
			return err, sterr.Alternate
		case 400:
			err = errors.New(errdef.ErrBase + "malformed STUN request: ")
			log.Println(err.Error(), declineStatus, ", reason: ", reason)
			return err, sterr.Other
		case 500:
			err = errors.New(errdef.ErrBase + "temporary server error: ")
			log.Println(err.Error(), declineStatus, ", reason: ", reason)
			return err, sterr.Retry
		default:
			err = errors.New(errdef.ErrBase + "unexpected error: ")
			log.Println(err.Error(), declineStatus, ", reason: ", reason)
			return err, sterr.Alternate
		}
	} else {
		err = errors.New(errdef.ErrBase + "unmatching STUN response type")
		log.Println(err.Error())
		return err, sterr.Other
	}
	return nil, sterr.Success
}
