package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"fmt"

	"code.google.com/p/go.crypto/pbkdf2"
)

type AES struct {
	salt, key []byte
}

// AES encryption
func NewAES(password string) (*AES, error) {
	passwd := []byte(password)
	salt := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	key := pbkdf2.Key(passwd, salt, 4096, 32, sha1.New)
	fmt.Printf("KEY  : %x [%d]\n", key, len(key))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return new(AES), nil
}

func (a *AES) Encrypt(plaintext []byte) []byte {
	plain := make([]byte, 0)
	if len(plaintext)%aes.BlockSize != 0 {
		padding = secret.PKCS7Padding(len(plaintext), aes.BlockSize)
		plain := make([]byte, len(plaintext)+len(padding))
		copy(plain, plaintext)
		copy(plain[len(plaintext):], padding)
	}
	cipher := make([]byte, aes.BlockSize+len(plain))
	return cipher
}

func (a *AES) Decrypt(ciphertext []byte) []byte {
	return nil
}
