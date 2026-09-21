package custom

import (
	"errors"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/mikonova/OpenGate/errdef"
	"github.com/mikonova/OpenGate/protocol/custom/async"
	"github.com/mikonova/OpenGate/protocol/custom/payloads"
	"github.com/mikonova/OpenGate/protocol/custom/payloads/types"
	"github.com/mikonova/OpenGate/protocol/stun"
	"github.com/mikonova/OpenGate/ticker"
)

var ClientHello payloads.Packet = payloads.Packet{
	PacketType: types.Hello,
}

var KeepAlive payloads.Packet = payloads.Packet{
	PacketType: types.KeepAlive,
}

// control Channel
var connections map[string]chan int

func GlobalConnectionsInit() {
	connections = make(map[string]chan int)
}

// send a signal to exit a goroutine and close an async connection
func CloseConn(addrInfo stun.AddrInfo) {
	connections[addrInfo.Address.String()] <- -1
}

// establish a new connection and delegate a new goroutine which sends keepalive packets
func InitializeConn(addrInfo stun.AddrInfo, pt *ticker.PrecisionTicker) *async.SafeConn {
	var (
		safeConn *async.SafeConn
		err      error
	)
	connections[addrInfo.Address.String()] = make(chan int, 1)

	for tick := range pt.Tick {
		if tick.Second() == 0 { // sync with the new second
			safeConn, err = callClient(addrInfo)
			if err != nil {
				log.Println(errdef.WarnBase, err.Error())
				continue
			} else {
				break
			}
		}
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		defer safeConn.Close()
		for range ticker.C {
			if n := <-connections[addrInfo.Address.String()]; n == -1 {
				break
			}
			if _, err := safeConn.Write(KeepAlive.Encode()); err != nil {
				log.Println(errdef.ErrBase, err.Error())
			}
		}
	}()
	return safeConn
}

func callClient(aInfo stun.AddrInfo) (*async.SafeConn, error) {
	conn, err := net.Dial("udp", aInfo.Address.String()+strconv.Itoa(int(aInfo.Port)))
	if err != nil {
		return &async.SafeConn{}, errors.New(errdef.WarnBase + "client connection error")
	}
	return async.MakeSafeConn(conn), nil
}
