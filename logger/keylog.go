package logger

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"

	"github.com/mikonova/OpenGate/errdef"
	"github.com/mikonova/goenv"
)

func KeyGen(appPath string) (path string) {
	path = appPath + ".env"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC|os.O_APPEND, 0640)
	if err != nil {
		log.Fatalln(errdef.ErrBase, err)
	}
	defer file.Close()
	curve := ecdh.P256()
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		log.Println(errdef.ErrBase, "cannot create private key")
		defer KeyGen(appPath)
	}

	publicKey := privateKey.PublicKey()
	privKeyHex := hex.EncodeToString(privateKey.Bytes())
	pubKeyHex := hex.EncodeToString(publicKey.Bytes())
	file.WriteString("public=" + pubKeyHex + "\n" + "secret=" + privKeyHex + "\n")
	return
}

func MoveKeyToEnv(path string) {
	goenv.FetchFiles(path)
	goenv.ParseFiles()
}
