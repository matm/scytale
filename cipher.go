package secret

type Cipher interface {
	Encrypt(key, clear []byte) []byte
	Decrypt(key, clear []byte) []byte
}
