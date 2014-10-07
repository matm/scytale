package main

import (
	"fmt"
	"os"
	"strings"

	"secret"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s key plain\n", os.Args[0])
		os.Exit(2)
	}
	password := strings.ToUpper(os.Args[1])
	clear := strings.ToUpper(os.Args[2])

	v := secret.NewVigenere()

	fmt.Println("CLEAR ", clear)

	key := make([]byte, len(clear))
	if len(password) < len(clear) {
		j := 0
		for k := 0; k < len(clear); k++ {
			key[k] = password[j]
			j += 1
			if j > len(password)-1 {
				j = 0
			}
		}
	}
	fmt.Println("KEY   ", string(key))

	cipher := v.Encrypt([]byte(key), []byte(clear))
	fmt.Println("CIPHER", string(cipher))

	plain := v.Decrypt([]byte(key), cipher)
	fmt.Println("CLEAR ", string(plain))

	fmt.Print("\nXOR cipher\n")
	x := secret.NewXor()
	a := []byte{1, 2, 3, 5, 6, 7}
	k := []byte{212, 16, 24, 32, 68, 44}
	fmt.Println("CLEAR ", a)
	fmt.Println("KEY   ", k)
	z := x.Encrypt(k, a)
	fmt.Println("CIPHER", z)
	fmt.Println("CLEAR ", x.Decrypt(k, z))
}
