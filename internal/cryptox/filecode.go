package cryptox

import (
	"crypto/rand"
	"math/big"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
const codeLength = 6

func GenerateFileCode() (string, error) {
	buffer := make([]byte, codeLength)

	charsetLen := big.NewInt(int64(len(charset)))
	for i := 0; i < codeLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		buffer[i] = charset[randomIndex.Int64()]
	}

	return string(buffer), nil
}
