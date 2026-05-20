package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	PasswordHash string `json:"password_hash"`
}

var (
	DiaryDir   string
	ConfigFile string
)

func InitPaths() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	DiaryDir = filepath.Join(home, "Documents", "virt-diary")
	ConfigFile = filepath.Join(DiaryDir, ".config.json")
	return os.MkdirAll(DiaryDir, 0700)
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile, data, 0600)
}
