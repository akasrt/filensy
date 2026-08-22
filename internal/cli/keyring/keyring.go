package keyring

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "filensy"
	accountName = "localstore-password"
)

func GetLocalPassword() (string, error) {
	password, err := keyring.Get(serviceName, accountName)
	if err != nil {
		return "", fmt.Errorf("failed to get password: %w", err)
	}
	return password, nil
}

func GenerateAndSetPassword() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	newPassword := base64.StdEncoding.EncodeToString(randomBytes)

	err := keyring.Set(serviceName, accountName, newPassword)
	if err != nil {
		return "", fmt.Errorf("failed to save password to keyring: %w", err)
	}

	return newPassword, nil
}
