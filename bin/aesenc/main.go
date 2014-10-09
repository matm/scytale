package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"secret"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -o output input password\n", os.Args[0])
		flag.PrintDefaults()
	}
	help := flag.Bool("h", false, "show help message")
	decrypt := flag.Bool("d", false, "decrypt file")
	output := flag.String("o", "", "output file name")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(2)
	}
	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(2)
	}
	if *output == "" {
		log.Fatal("missing output file name (use -o)")
	}
	name := flag.Arg(0)
	fmt.Println("PWD", flag.Arg(1))
	a, err := secret.NewAES(flag.Arg(1))
	if err != nil {
		log.Fatal("AES init:", err)
	}

	var action func(name, output string) error
	action = a.EncryptFile
	if *decrypt {
		fmt.Println("DEC")
		action = a.DecryptFile
	}

	if err := action(name, *output); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote to", *output)
}
