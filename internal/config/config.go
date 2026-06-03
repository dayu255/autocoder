package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultContest  string     `json:"defaultContest"`
	DefaultLanguage string     `json:"defaultLanguage"`
	TemplateFile    []Template `json:"templateFile"`
}

type Template struct {
	Language string `json:"language"`
	FilePath string `json:"filePath"`
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// config.jsonの場所は~/.config/autocoder/config.json
	configPath := filepath.Join(home, ".config", "autocoder", "autocoder.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.New("config.json is not exist")
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, errors.New("Fail to open config file")
	}
	configData := make([]byte, 1024)
	_, err = configFile.Read(configData)
	if err != nil {
		return nil, errors.New("Fail to read config file")
	}

	var cfg Config
	json.Unmarshal(configData, &cfg)

	return &cfg, nil
}
