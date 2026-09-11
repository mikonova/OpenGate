package secure

import (
	_ "crypto/aes"
	_ "crypto/cipher"
	_ "crypto/ecdh"
	_ "crypto/rand"
	_ "log"
	"os"
)

func GetKeys() (secret, public string) {
	os.Getenv("secret")
	os.Getenv("public")
}

func CreateSharedSecret(remotePublic []byte, secret []byte) {

}

// func ProcessBytes(compositeKey []byte, data []byte) []byte {
// 	key := compositeKey[:32]
// 	iv := compositeKey[32:]

// 	dst := make([]byte, len(data))

// 	cip, err := aes.NewCipher(key)
// 	if err != nil {
// 		log.Panicln(errdef.ErrBase + err.Error())
// 	}
// 	stream := cipher.NewCTR(cip, iv)
// 	stream.XORKeyStream(dst, data)
// 	return dst
// }
