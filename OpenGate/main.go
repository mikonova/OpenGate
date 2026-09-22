package main

import (
	"fmt"
	"time"

	"github.com/mikonova/OpenGate/logger"
	"github.com/mikonova/OpenGate/protocol/custom"
	"github.com/mikonova/OpenGate/protocol/ntp"
)

var AppPath string

func main() {

	custom.GlobalConnectionsInit()
	AppPath = logger.SetDefDir()

	ticker := ntp.GetNtpTime(time.Second)
	fmt.Println("Тикер получен")
	time.Sleep(time.Second * 5)
	for tick := range ticker.Tick {
		fmt.Println(tick)
	}

}
