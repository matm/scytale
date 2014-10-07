package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
	"log"

	"code.google.com/p/go.crypto/pbkdf2"
	"secret"
)

func main() {
	passwd := []byte("mypasswd")
	salt := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	key := pbkdf2.Key(passwd, salt, 4096, 32, sha1.New)
	fmt.Printf("KEY  : %x [%d]\n", key, len(key))

	plaintext := []byte("This is a crazy message, man")

	padding := secret.PKCS7Padding(len(plaintext), aes.BlockSize)
	plain := make([]byte, len(plaintext)+len(padding))
	copy(plain, plaintext)
	copy(plain[len(plaintext):], padding)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plain))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}

	// Operation mode
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plain)

	fmt.Printf("PLAIN: %v [%d]\n", string(plaintext), len(plaintext))
	fmt.Printf("CRYPT: %x [%d]\n", ciphertext, len(ciphertext))

	// Decrypt
	mode = cipher.NewCBCDecrypter(block, ciphertext[:aes.BlockSize])
	ciphertext = ciphertext[aes.BlockSize:]
	mode.CryptBlocks(ciphertext, ciphertext)

	// Remove padding
	cnt := ciphertext[len(ciphertext)-1]
	ciphertext = ciphertext[:len(ciphertext)-int(cnt)]
	fmt.Printf("PLAIN: %v [%d]\n", string(ciphertext), len(ciphertext))
}
