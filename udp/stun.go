package udp

import (
	"crypto/rand"
	"log"
	"net"
	"time"

	"github.com/mikonova/OpenGate/errdef"
	"github.com/mikonova/OpenGate/udp/stunlist"
)

type addrInfo struct {
	IsIPv6  bool
	Port    string
	Address string
}

func StunDial() {
	var conn net.Conn
	transactionID := make([]byte, 12)
	inputBuf := make([]byte, 0)
	buf := []byte{
		0x00, 0x01, // binding request
		0x21, 0x12, 0xA4, 0x42, // magic cookie (stun)
	}

	_, err := rand.Read(transactionID)
	if err != nil {
		log.Fatalln(err)
	}

	buf = append(buf, transactionID...)

	for _, v := range stunlist.ServerList {
		conn, err = net.Dial("udp", v)
		if err != nil {
			conn.Close()
			log.Println(errdef.WarnBase + "dialing error on server \"" + v + "\", switching server")
			continue
		}
		conn.SetDeadline(time.Now().Add(time.Second * 5))
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
		decodeResp(inputBuf)

	}

}

func decodeResp(buf []byte) {
	isIPv6 := false

	status := buf[:1]
	family := buf[1:2]
	xPort := buf[2:4]
	xAttribs := buf[4:]
	if family[0] == byte(0x01) {
		isIPv6 = true
	}
	
	if xPort[0] == 

}
