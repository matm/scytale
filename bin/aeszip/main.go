package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"log"
	"os"

	"secret"
)

const pwdMinLen = 4

func walk(path string, info os.FileInfo, current, total int) error {
	fmt.Printf("[%02d/%02d] %-70s\r", current, total, info.Name())
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s
  [[-p password] -o output.zip filepattern]
  [-x [-p password] [-o output_dir] [-n pos] archive.zip]
  [-l archive.zip]
  [-s archive.zip]

where options are
`, os.Args[0])
		flag.PrintDefaults()
	}
	help := flag.Bool("h", false, "show help message")
	output := flag.String("o", "", "output zip archive file")
	extract := flag.Bool("x", false, "extract and decrypt files")
	list := flag.Bool("l", false, "list files in archive")
	stats := flag.Bool("s", false, "archive stats")
	pos := flag.Int("n", -1, "extract file at pos in archive")
	password := flag.String("p", "", "password to use (UNSECURE)")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(2)
	}
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if *stats {
		r, err := zip.OpenReader(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()
		fmt.Printf("%d\n", len(r.File))
		return
	}
	if *list {
		r, err := zip.OpenReader(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()

		fmt.Printf(`Archive:  %s
  Length      Date    Time    Name
---------  ---------- -----   ----
`, flag.Arg(0))
		var total int64
		total = 0
		for _, f := range r.File {
			info := f.FileInfo()
			fmt.Printf("%9d  %s   %s\n", info.Size(),
				info.ModTime().Format("2006-01-02 15:04"), f.Name)
			total += info.Size()
		}
		fmt.Println("---------  		      -------")
		fmt.Printf("%9d  		      %d files\n", total, len(r.File))
		return
	}
	if *extract {
		if *password == "" {
			pwd, err := secret.ReadPassword(pwdMinLen, false)
			if err != nil {
				log.Fatal(err)
			}
			*password = pwd
		}
		ar := secret.NewZipArchive(*password)
		if *output == "" {
			*output = "."
		}
		if *pos >= 0 {
			// Extract at pos only
			if err := ar.ExtractAt(*pos, flag.Arg(0), *output); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := ar.ExtractAll(flag.Arg(0), *output); err != nil {
				log.Fatal(err)
			}
		}
		return
	}
	if *output == "" {
		log.Fatal("missing output file name (use -o)")
	}

	if *password == "" {
		pwd, err := secret.ReadPassword(pwdMinLen, false)
		if err != nil {
			log.Fatal(err)
		}
		*password = pwd
	}
	ar := secret.NewZipArchive(*password)
	if err := ar.Create(*output, flag.Args(), walk); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nWrote to %s\n", *output)
}
