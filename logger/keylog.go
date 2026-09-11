package logger

import (
	cryptrand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"log"
	"math/big"
	"math/rand"
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
	random, err := cryptrand.Int(cryptrand.Reader, big.NewInt(1<<32))
	byteKey := make([]byte, 0)
	binary.LittleEndian.AppendUint64(byteKey, random.Uint64())
	keyHash := sha256.Sum256(byteKey)
	keySlice := keyHash[:]
	byteNonce := make([]byte, 0)
	binary.LittleEndian.AppendUint64(byteNonce, rand.Uint64())
	nonceArray := sha256.Sum256(byteNonce)
	hashedNonce := nonceArray[:]

	compositeKey := hex.EncodeToString(keySlice) + ";" + hex.EncodeToString(hashedNonce)
	file.WriteString("key=" + compositeKey)
	return
}

func MoveKeyToEnv(path string) {
	goenv.FetchFiles(path)
	goenv.ParseFiles()
}
