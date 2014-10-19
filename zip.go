package secret

import (
	"archive/zip"
	"os"
	"path"
	"sync"
)

type ZipArchive struct {
	sync.Mutex
	password  string
	processed []string
	fileset   []string
}

func NewZipArchive(password string) *ZipArchive {
	return &ZipArchive{password: password}
}

func (a *ZipArchive) SetPassword(pwd string) {
	a.password = pwd
}

func (a *ZipArchive) Create(output string, files []string) error {
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	tw := zip.NewWriter(out)
	defer tw.Close()

	aes, err := NewAES(a.password)
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
			f.Close()
			return err
		}
		hdr, err := zip.FileInfoHeader(info)
		fw, err := tw.CreateHeader(hdr)
		if err != nil {
			f.Close()
			return err
		}
		if err := aes.EncryptFile(f, fw); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}

func (a *ZipArchive) Extract(input, outputDir string) error {
	if _, err := os.Stat(outputDir); err != nil {
		return err
	}
	tr, err := zip.OpenReader(input)
	if err != nil {
		return err
	}
	defer tr.Close()

	aes, err := NewAES(a.password)
	if err != nil {
		return err
	}

	for _, fi := range tr.File {
		f, err := os.Create(path.Join(outputDir, fi.Name))
		if err != nil {
			f.Close()
			return err
		}
		rc, err := fi.Open()
		if err != nil {
			return err
		}
		if err := aes.DecryptFile(rc, f); err != nil {
			f.Close()
			return err
		}
		f.Close()
		rc.Close()
	}
	return nil
}
