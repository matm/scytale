package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/matm/scytale"
)

const pwdMinLen = 4

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
		password, err := scytale.ReadPassword(pwdMinLen, false)
		if err != nil {
			log.Fatal(err)
		}
		ar := scytale.NewTarArchive(password)
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

	password, err := scytale.ReadPassword(pwdMinLen, true)
	if err != nil {
		log.Fatal(err)
	}
	ar := scytale.NewTarArchive(password)
	if err := ar.Create(*output, flag.Args()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote to", *output)
}
