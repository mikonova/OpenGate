package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/mikonova/OpenGate/logger"
	"github.com/mikonova/OpenGate/protocol/ntp"
)

var AppPath string

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	AppPath = logger.SetDefDir()

	ticker := ntp.GetNtpTime(time.Second)
	fmt.Println("Тикер получен")
	time.Sleep(time.Second * 5)
	for tick := range ticker.Tick {
		fmt.Println(tick)
	}

}
