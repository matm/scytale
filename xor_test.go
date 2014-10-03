package secret

import (
	"testing"
)

func TestXor(t *testing.T) {
	batch := []struct {
		clear, key, cipher []byte
	}{
		{[]byte{4}, []byte{0}, []byte{4}},
		{[]byte{2}, []byte{2}, []byte{0}},
		{
			[]byte{1, 2, 3, 5, 6, 7},
			[]byte{212, 16, 24, 32, 68, 44},
			[]byte{213, 18, 27, 37, 66, 43},
		},
	}
	x := NewXor()
	for _, b := range batch {
		c := x.Encrypt(b.key, b.clear)
		for k := range b.cipher {
			if c[k] != b.cipher[k] {
				t.Errorf("expect byte %v at pos %d, got %v", b.cipher[k], k, c[k])
			}
		}
	}
}
