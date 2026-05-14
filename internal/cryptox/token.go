package cryptox

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
)

const tokenSize = 32

func GenerateFileToken() (token string, hash string, err error) {
	b := make([]byte, tokenSize)

	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}

	token = base64.RawURLEncoding.EncodeToString(b)

	sum := sha256.Sum256([]byte(token))
	hash = hex.EncodeToString(sum[:])

	return token, hash, nil
}

func VerifyToken(token, hash string) bool {
	sum := sha256.Sum256([]byte(token))
	computedHash := hex.EncodeToString(sum[:])

	return subtle.ConstantTimeCompare([]byte(computedHash), []byte(hash)) == 1
}
