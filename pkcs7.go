package scytale

// PKCS#7 data padding
func PKCS7Padding(clearTextLen, blockSize int) []byte {
	numBytes := blockSize - (clearTextLen % blockSize)
	padding := make([]byte, numBytes)
	for j := 0; j < numBytes; j++ {
		padding[j] = byte(numBytes)
	}
	return padding
}
