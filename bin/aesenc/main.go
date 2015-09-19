package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/matm/scytale.v1"
)

const pwdMinLen = 4

// App version
const appVersion = "1.0"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -o output input\n", os.Args[0])
		flag.PrintDefaults()
	}
	help := flag.Bool("h", false, "show help message")
	decrypt := flag.Bool("d", false, "decrypt file")
	output := flag.String("o", "", "output file name")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(2)
	}
	if *version {
		fmt.Println("version:", appVersion)
		os.Exit(0)
	}
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(2)
	}
	if *output == "" {
		log.Fatal("missing output file name (use -o)")
	}
	// Read password twice
	twice := true
	if *decrypt {
		twice = false
	}
	pwd, err := scytale.ReadPassword(pwdMinLen, twice)
	if err != nil {
		log.Fatal(err)
	}
	name := flag.Arg(0)
	a, err := scytale.NewAES(pwd)
	if err != nil {
		log.Fatal("AES init:", err)
	}

	in, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()
	out, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	var action func(io.Reader, io.Writer) error
	action = a.EncryptFile
	if *decrypt {
		action = a.DecryptFile
	}

	if err := action(in, out); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote to", *output)
}
