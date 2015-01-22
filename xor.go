package scytale

type Xor int

// Basic XOR cipher implementation.
func NewXor() Cipher {
	return new(Xor)
}

func (x *Xor) Encrypt(key, clear []byte) []byte {
	cipher := make([]byte, len(clear))
	for k := 0; k < len(clear); k++ {
		cipher[k] = clear[k] ^ key[k]
	}
	return cipher
}

func (x *Xor) Decrypt(key, cipher []byte) []byte {
	return x.Encrypt(key, cipher)
}
