package main

import (
	"archive/tar"
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
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(2)
	}
	/*
		if len(flag.Args()) != 2 {
			flag.Usage()
			os.Exit(2)
		}
		if *output == "" {
			log.Fatal("missing output file name (use -o)")
		}
		name := flag.Arg(0)
		a, err := secret.NewAES(flag.Arg(1))
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
	*/
	out, err := os.Create("out.tar")
	if err != nil {
		log.Fatal(err)
	}
	tw := tar.NewWriter(out)
	defer tw.Close()

	a, err := secret.NewAES("passwd")
	if err != nil {
		log.Fatal("AES init:", err)
	}

	for _, file := range flag.Args() {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		info, err := f.Stat()
		if err != nil {
			log.Fatal(err)
		}
		hdr, err := tar.FileInfoHeader(info, "")
		hdr.Size = hdr.Size + 32 + 16
		if err != nil {
			log.Fatal(err)
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal(err)
		}
		if err := a.EncryptFile(f, tw); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Wrote", info.Name())
	}

	/*
		// Open the tar archive for reading.
		r := bytes.NewReader(buf.Bytes())
		tr := tar.NewReader(r)

		// Iterate through the files in the archive.
		for {
			_, err := tr.Next()
			if err == io.EOF {
				// end of tar archive
				break
			}
			if err != nil {
				log.Fatalln(err)
			}
		}
	*/
}
