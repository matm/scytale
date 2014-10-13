package secret

import (
	"archive/tar"
	"io"
	"os"
	"path"
)

type Archive struct {
	password string
}

func NewArchive(password string) *Archive {
	return &Archive{password}
}

func (a *Archive) Create(output string, files []string) error {
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	tw := tar.NewWriter(out)
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
			return err
		}
		hdr, err := tar.FileInfoHeader(info, "")
		//hdr.Name = fmt.Sprintf("%04d.crypt", j+1)
		// Estimate length of encrypted file
		hdr.Size = aes.EncryptedFileLength(info)
		if err != nil {
			return err
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if err := aes.EncryptFile(f, tw); err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

func (a *Archive) Extract(input, outputDir string) error {
	if _, err := os.Stat(outputDir); err != nil {
		return err
	}
	src, err := os.Open(input)
	if err != nil {
		return err
	}
	tr := tar.NewReader(src)

	aes, err := NewAES(a.password)
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
		f, err := os.Create(path.Join(outputDir, hdr.Name))
		if err != nil {
			return err
		}
		if err := aes.DecryptFile(tr, f); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}
