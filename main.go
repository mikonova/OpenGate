package main

import (
	"fmt"

	"github.com/mikonova/OpenGate/logger"
	"github.com/mikonova/OpenGate/stun"
)

func main() {
	logger.SetDefDir()
	info := stun.StunDial()
	fmt.Println("address: ", info.Address, "port: ", info.Port)
}
