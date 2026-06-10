package userconfig

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

const (
	appName  = "filensy"
	fileName = "config.json"
)

type Config struct {
	Directory string `json:"directory"`
	AuthKey   string `json:"auth_key"`
}

var (
	mu      sync.Mutex
	cfg     Config
	cfgPath string
)

func initPath() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dir = filepath.Join(dir, appName)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	cfgPath = filepath.Join(dir, fileName)

	return nil
}

func Load() error {
	err := initPath()
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfg = Config{}
			return nil
		} else {
			return err
		}
	}

	if string(data) == "" {
		cfg = Config{}
		return nil
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}

	return nil
}

func GetConfig() Config {
	return cfg
}

func SetConfig(conf Config) error {
	mu.Lock()
	defer mu.Unlock()

	cfg = conf
	err := saveLocked()
	if err != nil {
		return err
	}

	return nil
}

func saveLocked() error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(cfgPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
