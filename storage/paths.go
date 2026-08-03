package storage

import (
	"os"
	"path/filepath"
)

// GetDbmwDir returns the path to ~/.dbmw directory, creating it if needed.
func GetDbmwDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".dbmw")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return dir, nil
}
