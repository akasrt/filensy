package userconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/akasrt/filensy/internal/cli/keyring"
	"github.com/akasrt/filensy/internal/cryptox"
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

	localPassword, err := keyring.GetLocalPassword()
	if err != nil {
		cfg.AuthKey = ""
	}

	if cfg.AuthKey != "" {
		decryptedKey, err := cryptox.DecryptValue(cfg.AuthKey, localPassword)
		if err != nil {
			cfg.AuthKey = ""
			return fmt.Errorf("failed to decrypt auth key: %w", err)
		}
		cfg.AuthKey = decryptedKey
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
	diskCfg := cfg

	if diskCfg.AuthKey != "" {
		localPassword, err := keyring.GetLocalPassword()
		if err != nil || localPassword == "" {
			return fmt.Errorf("unable to save config as local password isn't available! Err: %w", err)
		}

		encryptedKey, err := cryptox.EncryptValue(diskCfg.AuthKey, localPassword)
		if err != nil {
			return fmt.Errorf("failed to encrypt auth key! Err: %w", err)
		}

		diskCfg.AuthKey = encryptedKey
	}

	data, err := json.MarshalIndent(diskCfg, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(cfgPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
