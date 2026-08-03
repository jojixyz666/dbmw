package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"dbmw/core/connection"
)

// ConnectionStore implements connection.ConnectionRepository using ~/.dbmw/connections.json.
type ConnectionStore struct {
	filePath string
	mu       sync.RWMutex
}

// NewConnectionStore creates a connection repository file store.
func NewConnectionStore(customPath ...string) (*ConnectionStore, error) {
	var targetPath string
	if len(customPath) > 0 && customPath[0] != "" {
		targetPath = customPath[0]
	} else {
		dir, err := GetDbmwDir()
		if err != nil {
			return nil, err
		}
		targetPath = filepath.Join(dir, "connections.json")
	}

	return &ConnectionStore{filePath: targetPath}, nil
}

func (s *ConnectionStore) readAll() ([]connection.ConnectionConfig, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []connection.ConnectionConfig{}, nil
		}
		return nil, err
	}
	var items []connection.ConnectionConfig
	if err := json.Unmarshal(data, &items); err != nil {
		return []connection.ConnectionConfig{}, nil
	}
	return items, nil
}

func (s *ConnectionStore) writeAll(items []connection.ConnectionConfig) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0600)
}

// Get retrieves a connection profile by ID.
func (s *ConnectionStore) Get(id string) (*connection.ConnectionConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items, err := s.readAll()
	if err != nil {
		return nil, err
	}
	for _, it := range items {
		if it.ID == id {
			c := it
			return &c, nil
		}
	}
	return nil, nil
}

// GetAll returns all saved connection profiles.
func (s *ConnectionStore) GetAll() ([]connection.ConnectionConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readAll()
}

// Save saves or updates a connection profile.
func (s *ConnectionStore) Save(cfg connection.ConnectionConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	items, err := s.readAll()
	if err != nil {
		items = []connection.ConnectionConfig{}
	}

	found := false
	for i, it := range items {
		if it.ID == cfg.ID {
			items[i] = cfg
			found = true
			break
		}
	}
	if !found {
		items = append(items, cfg)
	}

	return s.writeAll(items)
}

// Delete removes a stored connection by ID.
func (s *ConnectionStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	items, err := s.readAll()
	if err != nil {
		return err
	}

	newItems := make([]connection.ConnectionConfig, 0, len(items))
	for _, it := range items {
		if it.ID != id {
			newItems = append(newItems, it)
		}
	}

	return s.writeAll(newItems)
}
