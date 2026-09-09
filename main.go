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

// время - разница с точным временем = тикер смещенный на точное время
// тикер включается в случайное время. Тикер должен смотреть время по оффсету : хх.х0 секунд
