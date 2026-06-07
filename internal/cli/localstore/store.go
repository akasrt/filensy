package localstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/akasrt/filensy/internal/util/errorx"
)

const (
	fileName = "local.json"
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
	err := l.saveLocked()
	if err != nil {
		return fmt.Errorf("%w: %w", errorx.ErrLocalDeletionFailed, err)
	}

	return nil
}

func (l *store) load() error {
	data, err := os.ReadFile(l.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			l.files = map[string]FileData{}
			return nil
		} else {
			return err
		}
	}

	err = json.Unmarshal(data, &l.files)
	if err != nil {
		return err
	}

	return nil
}

func (l *store) saveLocked() error {
	data, err := json.MarshalIndent(l.files, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(l.filePath, data, 0600)
	if err != nil {
		return err
	}

	return nil
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
