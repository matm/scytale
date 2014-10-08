package main

import (
	"fmt"
	"log"
	"os"

	"secret"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s filename password\n", os.Args[0])
		os.Exit(2)
	}
	name := os.Args[1]
	a, err := secret.NewAES(os.Args[2])
	if err != nil {
		log.Fatal("AES init:", err)
	}
	out := name + ".crypt"
	if err := a.EncryptFile(name, out); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote to", out)

	clear := name + ".clear"
	if err := a.DecryptFile(out, clear); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote to", clear)
}
