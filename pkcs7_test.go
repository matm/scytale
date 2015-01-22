package scytale

import (
	"crypto/aes"
	"crypto/des"
	"testing"
)

func TestPKCS7Padding(t *testing.T) {
	ins := []struct {
		length, blockSize int
		expectLen         byte
	}{
		{4227, aes.BlockSize, 13},
		{1612475, aes.BlockSize, 5},
		{272473, des.BlockSize, 7},
		{10, aes.BlockSize, 6},
	}
	for _, c := range ins {
		padding := PKCS7Padding(c.length, c.blockSize)
		if byte(len(padding)) != c.expectLen {
			t.Errorf("bad padding length. Expect %d, got %d", c.expectLen, len(padding))
		}
		for k := range padding {
			if padding[k] != c.expectLen {
				t.Errorf("at pos %d, expect byte %v, got %v", k, c.expectLen, padding[k])
			}
		}
	}
}
