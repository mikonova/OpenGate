package ntp

import (
	"log"
	"net"
	"time"

	"github.com/mikonova/OpenGate/errdef"
)

var ntpServers = []string{
	"time.google.com",
	"time1.google.com",
	"time2.google.com",
	"time3.google.com",
	"time4.google.com",
	"time.android.com",
	"time.cloudflare.com",
	"time.windows.com",
	"ntp0.ntp-servers.net",
	"ntp1.ntp-servers.net",
	"ntp2.ntp-servers.net",
	"ntp3.ntp-servers.net",
	"ntp4.ntp-servers.net",
	"ntp5.ntp-servers.net",
	"ntp6.ntp-servers.net",
	"ntp7.ntp-servers.net",
}
var req = make([]byte, 48)

func GetNtpTime() {
	req[0] = 0b00011011
	for {
		ntpConnLoop()
	}
}

func ntpConnLoop() []byte {

	buffer := make([]byte, 48)
	for _, v := range ntpServers {
		conn, err := net.Dial("udp", v)
		if err != nil {
			log.Println(errdef.WarnBase + "dialing error on ntp server \"" + v + "\", switching server")
			continue
		}
		conn.SetDeadline(time.Now().Add(time.Second * 3))

		if _, err = conn.Write(req); err != nil {
			log.Println(errdef.WarnBase + "writing timeout on ntp server \"" + v + "\", switching server")
			conn.Close()
			continue
		}
		if n, err := conn.Read(buffer); err != nil {
			log.Println(errdef.WarnBase + "reading timeout on ntp server \"" + v + "\", switching server")
			conn.Close()
			continue
		} else if n != 48 {
			log.Println(errdef.WarnBase + "response malformed, switching server")
			conn.Close()
			continue
		} else {
			return buffer
		}
	}

}
