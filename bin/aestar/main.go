package main

import (
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

func getPassword(twice bool) string {
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
	if !twice {
		return password
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
	return password
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-o output.tar filepattern][-x -o output.tar archive.tar]\n", os.Args[0])
		flag.PrintDefaults()
	}
	help := flag.Bool("h", false, "show help message")
	output := flag.String("o", "", "output tar archive file")
	extract := flag.Bool("x", false, "extract and decrypt files")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(2)
	}
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if *extract {
		password := getPassword(false)
		ar := secret.NewArchive(password)
		if *output == "" {
			*output = "."
		}
		if err := ar.Extract(flag.Arg(0), *output); err != nil {
			log.Fatal(err)
		}
		return
	}
	if *output == "" {
		log.Fatal("missing output file name (use -o)")
	}

	password := getPassword(true)
	ar := secret.NewArchive(password)
	if err := ar.Create(*output, flag.Args()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote to", *output)
}
