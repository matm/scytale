package scytale

type Cipher interface {
	Encrypt(key, clear []byte) []byte
	Decrypt(key, cipher []byte) []byte
}
