package main

import (
	"fmt"
	"os"
	"strings"

	"crypto"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s key plain\n", os.Args[0])
		os.Exit(2)
	}
	password := strings.ToUpper(os.Args[1])
	clear := strings.ToUpper(os.Args[2])
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

	cipher := encrypt([]byte(key), []byte(clear))
	fmt.Println("CIPHER", string(cipher))

	plain := decrypt([]byte(key), cipher)
	fmt.Println("CLEAR ", string(plain))
}
