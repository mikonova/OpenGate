package secure

import (
	_ "crypto"
	_ "crypto/aes"
	_ "crypto/cipher"
	_ "crypto/rand"
	_ "crypto/sha256"
	_ "os"
)

func ProcessKey(key string) {
	//sha256.Sum256(key[])
}

func Encrypt(key []byte) {
	//cip := aes.NewCipher(key)
}
