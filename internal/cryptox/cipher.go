package cryptox

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	keyLength    = 32

	saltSize  = 16
	nonceSize = 12
)

// TODO: add a streaming cipher for large files

func Encrypt(reader io.Reader, writer io.Writer, password string) error {
	plain, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	salt, err := generateSalt(saltSize)
	if err != nil {
		return err
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nil, nonce, plain, nil)

	if _, err := writer.Write(salt); err != nil {
		return err
	}
	if _, err := writer.Write(nonce); err != nil {
		return err
	}
	if _, err := writer.Write(ciphertext); err != nil {
		return err
	}

	return nil
}

func Decrypt(reader io.Reader, writer io.Writer, password string) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if len(data) < saltSize+nonceSize {
		return io.ErrUnexpectedEOF
	}

	salt := data[:saltSize]
	nonce := data[saltSize : saltSize+nonceSize]
	ciphertext := data[saltSize+nonceSize:]

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	_, err = writer.Write(plain)
	if err != nil {
		return err
	}

	return nil
}

func deriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, keyLength)
}

func generateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func EncryptValue(plaintext string, password string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	var buf bytes.Buffer
	err := Encrypt(strings.NewReader(plaintext), &buf, password)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DecryptValue(ciphertextBase64 string, password string) (string, error) {
	if ciphertextBase64 == "" {
		return "", nil
	}

	rawCiphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = Decrypt(bytes.NewReader(rawCiphertext), &buf, password)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
