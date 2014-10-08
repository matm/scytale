package main

import (
	"fmt"
	"log"

	"secret"
)

func main() {
	a, err := secret.NewAES("mypasswd")
	plaintext := []byte("This is a crazy message, man")
	if err != nil {
		log.Fatal(err)
	}
	// Encrypt
	iv, err := a.InitEncryption()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PLAIN: %v [%d]\n", string(plaintext), len(plaintext))
	ciphertext := a.Encrypt(plaintext)
	fmt.Printf("CRYPT: %x [%d]\n", ciphertext, len(ciphertext))

	// Decrypt
	a.InitDecryption(iv)
	clear := a.Decrypt(ciphertext)
	clear = a.RemovePadding(clear)
	fmt.Printf("PLAIN: %v [%d]\n", string(clear), len(clear))
}
