package relay

import (
	"encoding/binary"
	"net"
	"relayesp/protocol/iptransmit/types"
	"relayesp/wifi"
	"time"
)

type RelayTableRow struct {
	ClientOneConn, ClientTwoConn net.Conn
	ClientOneData, ClientTwoData string
}

type UnorderedClient struct {
	Client net.Conn
	Ipaddr string
}

func (rtw RelayTableRow) Sync() {

}

func (rtw RelayTableRow) Reject() {

}

func findClientPairs() {

}

func StoreClients() []net.Conn {
	clients := make([]UnorderedClient, 0)
	maxRetries := 0
	for conn := range wifi.ConnChan {
		conn.SetReadDeadline(time.Now().Add(time.Second * 30))
		var staticBuf = [512]byte{}
		buf := staticBuf[:]
	READ:
		_, err := conn.Read(buf)
		if err != nil && maxRetries <= 3 {
			println("Couldnt read from a client, retrying...")
			maxRetries++
			goto READ

		} else if err != nil && maxRetries > 3 {
			println("Couldnt read from a client, retries threshold reached")
			if resp := respondToClient(conn); !resp {
				println("skipping problematic client, client ip:", conn.LocalAddr().String())
				continue
			}
		}
		// orderedData := iptr.Packet{
		// 	MsgType: ,
		// }
	}

}

func respondToClient(conn net.Conn) (success bool) {
	outgoungBuf := make([]byte, 14)
	outgoungBuf[0] = types.ServerError
	outgoungBuf[9] = 1
	binary.BigEndian.Uint32(outgoungBuf[10:])
	conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
	_, err := conn.Write(outgoungBuf)
	if err != nil {
		println("unable to send a responce, skipping")
		return false
	}
	return true
}
