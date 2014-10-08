package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"

	"code.google.com/p/go.crypto/pbkdf2"
)

type AES struct {
	salt, key []byte
	block     cipher.Block
	mode      cipher.BlockMode
	iv        []byte
}

// AES encryption
func NewAES(password string) (*AES, error) {
	passwd := []byte(password)
	// TODO: random salt
	salt := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	key := pbkdf2.Key(passwd, salt, 4096, 32, sha1.New)
	fmt.Printf("KEY  : %x [%d]\n", key, len(key))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aes := &AES{salt: salt, key: key, block: block}

	return aes, nil
}

// Computes a random IV and set cipher's operation mode to CBC.
func (a *AES) InitEncryption() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	a.iv = iv
	a.mode = cipher.NewCBCEncrypter(a.block, iv)
	return iv, nil
}

func (a *AES) Encrypt(plaintext []byte) []byte {
	plain := plaintext
	if len(plaintext)%aes.BlockSize != 0 {
		padding := PKCS7Padding(len(plaintext), aes.BlockSize)
		plain = make([]byte, len(plaintext)+len(padding))
		copy(plain, plaintext)
		copy(plain[len(plaintext):], padding)
	}
	ciphertext := make([]byte, len(plain))
	a.mode.CryptBlocks(ciphertext, plain)
	return ciphertext
}

// Uses IV and set cipher's operation mode to CBC.
func (a *AES) InitDecryption(iv []byte) {
	a.iv = iv
	a.mode = cipher.NewCBCDecrypter(a.block, iv)
}

func (a *AES) Decrypt(ciphertext []byte) []byte {
	a.mode.CryptBlocks(ciphertext, ciphertext)
	return ciphertext
}

func (a *AES) RemovePadding(clear []byte) []byte {
	cnt := clear[len(clear)-1]
	clear = clear[:len(clear)-int(cnt)]
	return clear
}
