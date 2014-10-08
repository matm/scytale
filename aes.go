package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"os"

	"code.google.com/p/go.crypto/pbkdf2"
)

type AES struct {
	salt, key []byte
	block     cipher.Block
	mode      cipher.BlockMode
	iv        []byte
}

// 2K buffer
const bufLen = 2048

// Strong AES encryption, with a cipher operating in CBC mode,
// using a derived 256 bits key using PBKDF2.
func NewAES(password string) (*AES, error) {
	passwd := []byte(password)
	// Use a random 16 bytes salt
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := pbkdf2.Key(passwd, salt, 4096, 32, sha1.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aes := &AES{salt: salt, key: key, block: block}

	return aes, nil
}

// Computes a random IV and set cipher's operation mode to CBC.
func (a *AES) InitEncryption() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	a.iv = iv
	a.mode = cipher.NewCBCEncrypter(a.block, iv)
	return iv, nil
}

func (a *AES) Encrypt(plaintext []byte) []byte {
	plain := plaintext
	if len(plaintext)%aes.BlockSize != 0 {
		padding := PKCS7Padding(len(plaintext), aes.BlockSize)
		plain = make([]byte, len(plaintext)+len(padding))
		copy(plain, plaintext)
		copy(plain[len(plaintext):], padding)
	}
	ciphertext := make([]byte, len(plain))
	a.mode.CryptBlocks(ciphertext, plain)
	return ciphertext
}

// Uses IV and set cipher's operation mode to CBC.
func (a *AES) InitDecryption(iv []byte) {
	a.iv = iv
	a.mode = cipher.NewCBCDecrypter(a.block, iv)
}

// Decrypt decripts a block of ciphertext. You have to ensure
// this ciphertext is a multiple of AES's block size.
func (a *AES) Decrypt(ciphertext []byte) []byte {
	a.mode.CryptBlocks(ciphertext, ciphertext)
	return ciphertext
}

// RemovePadding removes extra padding for plain text.
func (a *AES) RemovePadding(clear []byte) []byte {
	cnt := clear[len(clear)-1]
	clear = clear[:len(clear)-int(cnt)]
	return clear
}

// EncryptFile encrypts infile and saves the resulting AES encoding
// to outfile.
func (a *AES) EncryptFile(infile, outfile string) error {
	f, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer f.Close()
	out, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer out.Close()

	buf := make([]byte, bufLen)

	// Encrypt
	iv, err := a.InitEncryption()
	if err != nil {
		return err
	}
	_, err = out.Write(iv)
	if err != nil {
		return err
	}
	for {
		n, err := f.Read(buf)
		if n == 0 && err == io.EOF {
			break
		}
		if n < len(buf) {
			buf = buf[:n]
		}
		_, err = out.Write(a.Encrypt(buf))
		if err != nil {
			return err
		}
	}
	return nil
}

// DecryptFile decrypts infile and saves the resulting AES decoding
// to outfile.
func (a *AES) DecryptFile(infile, outfile string) error {
	f, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer f.Close()
	out, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer out.Close()

	buf := make([]byte, bufLen)

	// Decrypt
	iv := make([]byte, aes.BlockSize)
	_, err = f.Read(iv)
	if err != nil {
		return errors.New(fmt.Sprintf("can't read IV: %s", err.Error()))
	}
	a.InitDecryption(iv)

	var clear []byte
	for {
		n, err := f.Read(buf)
		if n == 0 && err == io.EOF {
			break
		}
		if n < len(buf) {
			// Last block, remove extra padding
			clear = a.Decrypt(buf[:n])
			clear = a.RemovePadding(clear)
		} else {
			clear = a.Decrypt(buf)
		}
		_, err = out.Write(clear)
		if err != nil {
			return err
		}
	}

	return nil
}
