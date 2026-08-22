package localstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/akasrt/filensy/internal/cli/keyring"
	"github.com/akasrt/filensy/internal/cryptox"
	"github.com/akasrt/filensy/internal/util/errorx"
)

const (
	fileName = "local-store"
	appName  = "filensy"
)

type Store interface {
	Create(data FileData) error
	Get(code string) (FileData, bool)
	Delete(code string) error
}

type store struct {
	filePath string
	mu       sync.RWMutex
	password string
	files    map[string]FileData
}

func New() (Store, error) {
	filePath, err := getFilePath()
	if err != nil {
		return nil, err
	}

	store := store{
		filePath: filePath,
		mu:       sync.RWMutex{},
		files:    map[string]FileData{},
	}

	err = store.load()
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (l *store) Create(data FileData) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.files[data.Code] = data

	if !l.isActive() {
		return errorx.ErrInactiveLocalStore
	}

	err := l.saveLocked()
	if err != nil {
		return fmt.Errorf("%w: %w", errorx.ErrLocalCreationFailed, err)
	}

	return nil
}

func (l *store) Get(code string) (FileData, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	data, exists := l.files[code]
	if !exists {
		return FileData{}, false
	}

	return data, true
}

func (l *store) Delete(code string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.files, code)

	if !l.isActive() {
		return errorx.ErrInactiveLocalStore
	}

	err := l.saveLocked()
	if err != nil {
		return fmt.Errorf("%w: %w", errorx.ErrLocalDeletionFailed, err)
	}

	return nil
}

func (l *store) load() error {
	password, err := keyring.GetLocalPassword()
	if err != nil || password == "" {
		l.files = map[string]FileData{}

		newPassword, err := keyring.GenerateAndSetPassword()
		if err != nil {
			fmt.Print(err)
			return fmt.Errorf("error creating local-store password")
		}

		l.password = newPassword
		return nil
	}
	l.password = password

	fileData, err := os.ReadFile(l.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			l.files = map[string]FileData{}
			return nil
		} else {
			return err
		}
	}

	if len(fileData) == 0 {
		l.files = map[string]FileData{}
		return nil
	}

	var decryptedBuf bytes.Buffer
	err = cryptox.Decrypt(bytes.NewReader(fileData), &decryptedBuf, l.password)
	if err != nil {
		os.Rename(l.filePath, l.filePath+".bak")
		l.files = map[string]FileData{}
		return nil
	}

	err = json.Unmarshal(decryptedBuf.Bytes(), &l.files)
	if err != nil {
		return err
	}

	return nil
}

func (l *store) saveLocked() error {
	if !l.isActive() {
		return errorx.ErrInactiveLocalStore
	}

	data, err := json.MarshalIndent(l.files, "", "  ")
	if err != nil {
		return err
	}

	var encryptedBuf bytes.Buffer
	err = cryptox.Encrypt(bytes.NewReader(data), &encryptedBuf, l.password)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	err = os.WriteFile(l.filePath, encryptedBuf.Bytes(), 0600)
	if err != nil {
		return err
	}

	return nil
}

func (l *store) isActive() bool {
	return l.password != ""
}

func getFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir = filepath.Join(dir, appName)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	path := filepath.Join(dir, fileName)

	return path, nil
}
