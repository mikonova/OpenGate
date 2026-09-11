package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
	"os"

	"github.com/mikonova/OpenGate/errdef"
)

func GetKey() string {
	return os.Getenv("key")
}

func ProcessBytes(compositeKey []byte, data []byte) []byte {
	key := compositeKey[:32]
	iv := compositeKey[32:]
	dst := make([]byte, len(data))

	cip, err := aes.NewCipher(key)
	if err != nil {
		log.Panicln(errdef.ErrBase + err.Error())
	}
	stream := cipher.NewCTR(cip, iv)
	stream.XORKeyStream(dst, data)
	return dst
}
