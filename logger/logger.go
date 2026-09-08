package logger

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mikonova/OpenGate/errdef"
)

func SetDefDir() {
	homepath, err := os.UserHomeDir()
	if err != nil {
		log.Println(errdef.ErrBase + "unable to find a user home dir")
		return
	}
	p := filepath.Join(homepath, "OpenGate")
	if _, err := os.Stat(p); err != nil {
		err = os.Mkdir(p, 0775)
		if err != nil {
			log.Fatalln(errdef.ErrBase + "unable to create a file dir")
		}
	}
	setStandartLogger(p)
}

func setStandartLogger(appPath string) {
	logPath := filepath.Join(appPath, "log")
	var err error
	if _, err := os.Stat(logPath); err != nil {
		err = os.Mkdir(logPath, 0775)
		if err != nil {
			log.Println(errdef.ErrBase + "unable to create a log dir")
			return
		}
	}

	logFilePath := filepath.Join(logPath, "OpenGate.log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR|os.O_APPEND, 0640)
	if err != nil {
		log.Println(errdef.ErrBase + "unable to create a logfile")
		return
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ltime | log.Lshortfile)

}
