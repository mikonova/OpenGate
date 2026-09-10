package logger

import (
	"os"

	_ "github.com/mikonova/goenv"
)

func KeyGen(appPath string) {
	file, err := os.OpenFile(appPath+".env", os.O_CREATE|os.O_RDWR|os.O_TRUNC|os.O_APPEND, 0640)
}
