package main

import (
	"gorelayesp/wifi"
	"machine"
	"time"
)

var (
	key, pass string
	CommChan  chan byte = make(chan byte)
)

func main() {
	machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200})
	time.Sleep(time.Second * 5)
LOGIN:
	println("enter wifi SSID")
	ssid := getInput()
	println("enter wifi password")
	pass := getInput()
	println("current inputs are: ", ssid, pass)
	println("proceed (y/n)?[y]")
	conf := getInput()
	if conf == "n" {
		goto LOGIN
	} else {
		wifi.WifiConnect(ssid, pass)
	}
	for {

	}
}

func getInput() (output string) {
	for {
		data, err := machine.Serial.ReadByte()
		if err != nil {
			continue
		}
		if data == '\n' {
			break
		} else {
			output += string(data)
		}
	}
	return output
}

/* func receiveControl() {
	for {
		input, err := machine.Serial.ReadByte()
	}
}
*/
