package secret

import (
	"archive/zip"
	"errors"
	"fmt"
	"os"
	"path"
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

func (a *ZipArchive) Create(output string, files []string, fn WalkFunc) error {
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
