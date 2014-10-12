package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"log"
	"os"

	"code.google.com/p/go.crypto/ssh/terminal"
	"secret"
)

func perror(msg string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	os.Exit(1)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -o output.tar filepattern\n", os.Args[0])
		flag.PrintDefaults()
	}
	help := flag.Bool("h", false, "show help message")
	output := flag.String("o", "", "output tar archive file")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(2)
	}
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if *output == "" {
		log.Fatal("missing output file name (use -o)")
	}

	fmt.Printf("Password: ")
	pwd, err := terminal.ReadPassword(int(os.Stdout.Fd()))
	if err != nil {
		perror(err.Error())
	}
	password := string(pwd)
	fmt.Println()
	if password == "" {
		perror("Empty password not allowed.")
	}
	fmt.Printf("Repeat: ")
	pwd2, err := terminal.ReadPassword(int(os.Stdout.Fd()))
	if err != nil {
		perror(err.Error())
	}
	fmt.Println()
	confirm := string(pwd2)
	if password != confirm {
		perror("Passwords mismatch.")
	}

	out, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	tw := tar.NewWriter(out)
	defer tw.Close()

	a, err := secret.NewAES(password)
	if err != nil {
		log.Fatal("AES init:", err)
	}

	for _, file := range flag.Args() {
		fmt.Printf("Adding %s ... ", file)
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		info, err := f.Stat()
		if err != nil {
			log.Fatal(err)
		}
		hdr, err := tar.FileInfoHeader(info, "")
		//hdr.Name = fmt.Sprintf("%04d.crypt", j+1)
		// Estimate length of encrypted file
		hdr.Size = a.EncryptedFileLength(info)
		if err != nil {
			log.Fatal(err)
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		if err := a.EncryptFile(f, tw); err != nil {
			log.Fatal(err)
		}
		f.Close()
		fmt.Println("OK")
	}
	fmt.Println("Wrote to", *output)
}
