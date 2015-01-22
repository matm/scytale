package scytale

import (
	"archive/zip"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type ZipArchive struct {
	sync.Mutex
	password string
	status   Status
}

type Status int

const (
	Running Status = iota
	Idle
)

// WalkFunc is the type of the function called for each file added to the archive.
type WalkFunc func(path string, info os.FileInfo, current, total int) error

func NewZipArchive(password string) *ZipArchive {
	return &ZipArchive{
		password: password,
		status:   Idle,
	}
}

func (a *ZipArchive) SetPassword(pwd string) {
	a.password = pwd
}

func (a *ZipArchive) Status() Status {
	a.Lock()
	defer a.Unlock()
	return a.status
}

// Cancel stops the current processing, if any.
func (a *ZipArchive) Cancel() {
	a.Lock()
	defer a.Unlock()
	a.status = Idle
}

// Create creates a new ZIP output archive given a list of input files. fn() is
// called just after a file has been encrypted and added to the archive. Set random
// to true to rename file with their md5 checksum.
func (a *ZipArchive) Create(output string, files []string, fn WalkFunc, random bool) error {
	a.Lock()
	a.status = Running
	a.Unlock()

	defer func() {
		a.Lock()
		a.status = Idle
		a.Unlock()
	}()

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

	nb := len(files)
	for u, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		info, err := f.Stat()
		if err != nil {
			f.Close()
			return err
		}
		info.IsDir()
		hdr, err := zip.FileInfoHeader(info)
		if random {
			m := md5.New()
			io.WriteString(m, hdr.Name)
			hdr.Name = fmt.Sprintf("%x%s", m.Sum(nil), filepath.Ext(hdr.Name))
		}
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
		if err := fn(file, info, u+1, nb); err != nil {
			return err
		}
		if a.status == Idle {
			// Processing has been cancelled
			break
		}
	}
	return nil
}

// ExtractAll extracts all encrypted files from zip archive. The
// resulting files are decrypted using the provided password.
func (a *ZipArchive) ExtractAll(archive, outputDir string) error {
	a.Lock()
	a.status = Running
	a.Unlock()

	defer func() {
		a.Lock()
		a.status = Idle
		a.Unlock()
	}()

	if _, err := os.Stat(outputDir); err != nil {
		return err
	}
	tr, err := zip.OpenReader(archive)
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
		if a.status == Idle {
			// Processing has been cancelled
			break
		}
	}
	return nil
}

// ExtractAt extracts a single file at position pos from zip archive. The
// resulting file is decrypted using the provided password.
func (a *ZipArchive) ExtractAt(pos int, archive, outputDir string) error {
	a.Lock()
	a.status = Running
	a.Unlock()

	defer func() {
		a.Lock()
		a.status = Idle
		a.Unlock()
	}()

	if _, err := os.Stat(outputDir); err != nil {
		return err
	}
	tr, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer tr.Close()

	if pos > len(tr.File)-1 || pos < 0 {
		return errors.New(fmt.Sprintf("position %d out of bounds (%d files in archive)", pos, len(tr.File)))
	}

	aes, err := NewAES(a.password)
	if err != nil {
		return err
	}

	at := tr.File[pos]
	f, err := os.Create(path.Join(outputDir, at.Name))

	if err != nil {
		f.Close()
		return err
	}
	rc, err := at.Open()
	if err != nil {
		return err
	}
	if err := aes.DecryptFile(rc, f); err != nil {
		f.Close()
		return err
	}
	f.Close()
	rc.Close()

	return nil
}
