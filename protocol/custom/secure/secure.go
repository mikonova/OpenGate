package secure

import (
	_ "crypto/aes"
	_ "crypto/cipher"
	"crypto/ecdh"
	_ "crypto/rand"
	_ "log"
	"os"
)

func GetKeys() (secret, public string) {
	secret = os.Getenv("secret")
	public = os.Getenv("public")
	return
}

func CreateSharedSecret(remotePublic []byte, secret []byte) {

}

func 

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
