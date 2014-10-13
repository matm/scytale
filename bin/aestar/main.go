package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"io"
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

func createArchive(output, password string, files []string) error {
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	tw := tar.NewWriter(out)
	defer tw.Close()

	a, err := secret.NewAES(password)
	if err != nil {
		return err
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		info, err := f.Stat()
		if err != nil {
			return err
		}
		hdr, err := tar.FileInfoHeader(info, "")
		//hdr.Name = fmt.Sprintf("%04d.crypt", j+1)
		// Estimate length of encrypted file
		hdr.Size = a.EncryptedFileLength(info)
		if err != nil {
			return err
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if err := a.EncryptFile(f, tw); err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

func extractArchive(input, password string) error {
	src, err := os.Open(input)
	if err != nil {
		return err
	}
	tr := tar.NewReader(src)

	a, err := secret.NewAES(password)
	if err != nil {
		return err
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return err
		}
		f, err := os.Create(hdr.Name)
		if err != nil {
			return err
		}
		if err := a.DecryptFile(tr, f); err != nil {
			f.Close()
			return err
		}
		f.Close()
		fmt.Println("extracted", hdr.Name)
	}
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-o output.tar filepattern][-x archive.tar]\n", os.Args[0])
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
		if err := extractArchive(flag.Arg(0), password); err != nil {
			log.Fatal(err)
		}
		return
	}
	if *output == "" {
		log.Fatal("missing output file name (use -o)")
	}

	password := getPassword(true)
	if err := createArchive(*output, password, flag.Args()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Wrote to", *output)
}
