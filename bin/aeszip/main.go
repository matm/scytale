package main

import (
	"archive/zip"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/matm/scytale"
)

const pwdMinLen = 4

// App version
const VERSION = "1.2"

func walk(path string, info os.FileInfo, current, total int) error {
	fmt.Printf("[%02d/%02d] %-70s\r", current, total, info.Name())
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s
  [[-r] [-p password] -o output.zip filepattern]
  [-x [-p password] [-o output_dir] [-n pos] archive.zip]
  [[-j] -l archive.zip]
  [[-j] -s [-n pos] archive.zip]

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
	version := flag.Bool("v", false, "show version")
	jsonFormat := flag.Bool("j", false, "JSON output")
	random := flag.Bool("r", false, "rename files with their md5 checksum")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(2)
	}
	if *version {
		fmt.Println(VERSION)
		return
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
		if *pos >= 0 {
			if *pos >= len(r.File) {
				log.Fatal("position out of bounds")
			}
			f := r.File[*pos]
			if *jsonFormat {
				info := &struct {
					Name string `json:"name"`
					Size int64  `json:"size"`
				}{f.FileInfo().Name(), f.FileInfo().Size()}
				d, err := json.Marshal(info)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(d))
			} else {
				fmt.Println(f.FileInfo().Size(), f.FileInfo().Name())
			}
			return
		}
		fmt.Printf("%d\n", len(r.File))
		return
	}
	if *list {
		r, err := zip.OpenReader(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()

		if *jsonFormat {
			type data struct {
				Name string `json:"name"`
				Size int64  `json:"size"`
				Date string `json:"date"`
			}
			p := make([]data, 0)
			for _, f := range r.File {
				info := f.FileInfo()
				p = append(p, data{f.Name, info.Size(), info.ModTime().Format("2006-01-02 15:04")})
			}
			d, err := json.Marshal(p)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(d))
		} else {

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
		}
		return
	}
	if *extract {
		if *password == "" {
			pwd, err := scytale.ReadPassword(pwdMinLen, false)
			if err != nil {
				log.Fatal(err)
			}
			*password = pwd
		}
		ar := scytale.NewZipArchive(*password)
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
		pwd, err := scytale.ReadPassword(pwdMinLen, false)
		if err != nil {
			log.Fatal(err)
		}
		*password = pwd
	}
	ar := scytale.NewZipArchive(*password)
	if err := ar.Create(*output, flag.Args(), walk, *random); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nWrote to %s\n", *output)
}
