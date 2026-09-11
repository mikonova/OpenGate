package custom

import (
	"errors"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/mikonova/OpenGate/errdef"
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

func InitializeConn(addrInfo stun.AddrInfo, pt *ticker.PrecisionTicker) {

	for tick := range pt.Tick {
		if tick.Second() == 0 { // синхронизация с новой секундой
			conn, err := callClient(addrInfo)
			if err != nil {
				go log.Println(err.Error())
				continue
			}
		}
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		for t := range ticker.C {

		}
	}()
}

func callClient(aInfo stun.AddrInfo) (*net.Conn, error) {
	conn, err := net.Dial("udp", aInfo.Address.String()+strconv.Itoa(int(aInfo.Port)))
	if err != nil {
		return &conn, errors.New(errdef.WarnBase + "client connection error")
	}
	return &conn, nil
}
