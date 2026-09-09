package main

import (
	"fmt"
	"time"

	"github.com/mikonova/OpenGate/logger"
	"github.com/mikonova/OpenGate/protocol/ntp"
)

func main() {
	logger.SetDefDir()

	ticker := ntp.GetNtpTime(time.Second)
	for tick := range ticker.Tick {
		fmt.Println(tick)
	}

}
