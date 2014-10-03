package secret

type Vigenere int

// Vigenere's cipher implementation.
func NewVigenere() Cipher {
	return new(Vigenere)
}

func (v *Vigenere) Encrypt(key, clear []byte) []byte {
	cipher := make([]byte, len(clear))
	for k := 0; k < len(clear); k++ {
		cipher[k] = 65 + (clear[k]+key[k])%26
	}
	return cipher
}

func (v *Vigenere) Decrypt(key, cipher []byte) []byte {
	clear := make([]byte, len(cipher))
	for k := 0; k < len(cipher); k++ {
		clear[k] = 65 + (26+(cipher[k]-key[k]))%26
	}
	return clear
}
