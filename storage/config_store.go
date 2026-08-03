package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// AppConfig holds global application preferences.
type AppConfig struct {
	Theme           string `json:"theme"`
	DefaultPageSize int    `json:"defaultPageSize"`
	ServerPort      int    `json:"serverPort"`
	AutoOpenBrowser bool   `json:"autoOpenBrowser"`
}

// ConfigStore handles reading and writing ~/.dbmw/config.json.
type ConfigStore struct {
	filePath string
	mu       sync.RWMutex
}

// NewConfigStore creates a new config store.
func NewConfigStore(customPath ...string) (*ConfigStore, error) {
	var targetPath string
	if len(customPath) > 0 && customPath[0] != "" {
		targetPath = customPath[0]
	} else {
		dir, err := GetDbmwDir()
		if err != nil {
			return nil, err
		}
		targetPath = filepath.Join(dir, "config.json")
	}

	return &ConfigStore{filePath: targetPath}, nil
}

// DefaultConfig returns default initial application configuration.
func DefaultConfig() AppConfig {
	return AppConfig{
		Theme:           "dark",
		DefaultPageSize: 25,
		ServerPort:      8085,
		AutoOpenBrowser: true,
	}
}

// Get returns the loaded configuration or default if file doesn't exist.
func (s *ConfigStore) Get() (AppConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return DefaultConfig(), err
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}
	return cfg, nil
}

// Save writes app preferences to disk.
func (s *ConfigStore) Save(cfg AppConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0600)
}
