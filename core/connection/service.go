package connection

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Service provides application operations for database connections.
type Service struct {
	repo       ConnectionRepository
	factory    ConnectorFactory
	mu         sync.RWMutex
	connectors map[string]Connector
	configs    map[string]ConnectionConfig
	activeID   string
}

// NewService creates a new connection service.
func NewService(repo ConnectionRepository, factory ConnectorFactory) *Service {
	return &Service{
		repo:       repo,
		factory:    factory,
		connectors: make(map[string]Connector),
		configs:    make(map[string]ConnectionConfig),
	}
}

// TestConnection tests whether a connection can successfully be established.
func (s *Service) TestConnection(ctx context.Context, cfg ConnectionConfig) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	conn, err := s.factory(cfg.Driver)
	if err != nil {
		return fmt.Errorf("factory failed: %w", err)
	}
	defer conn.Close()

	testCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := conn.Connect(testCtx, cfg); err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}
	if err := conn.Ping(testCtx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}

// SaveConnection validates and saves a connection configuration.
func (s *Service) SaveConnection(cfg ConnectionConfig) (*ConnectionConfig, error) {
	if cfg.ID == "" {
		cfg.ID = fmt.Sprintf("conn_%d", time.Now().UnixNano())
		cfg.CreatedAt = time.Now()
	}
	cfg.UpdatedAt = time.Now()

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Save(cfg); err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.configs[cfg.ID] = cfg
	// Invalidate cached connector if config changed
	if old, exists := s.connectors[cfg.ID]; exists {
		old.Close()
		delete(s.connectors, cfg.ID)
	}
	s.mu.Unlock()

	return &cfg, nil
}

// GetConnection returns a saved connection by ID.
func (s *Service) GetConnection(id string) (*ConnectionConfig, error) {
	s.mu.RLock()
	c, ok := s.configs[id]
	s.mu.RUnlock()
	if ok {
		return &c, nil
	}
	return s.repo.Get(id)
}

// ListConnections returns all saved connections.
func (s *Service) ListConnections() ([]ConnectionConfig, error) {
	return s.repo.GetAll()
}

// DeleteConnection removes a saved connection and closes its active connector if open.
func (s *Service) DeleteConnection(id string) error {
	s.mu.Lock()
	if conn, exists := s.connectors[id]; exists {
		conn.Close()
		delete(s.connectors, id)
	}
	delete(s.configs, id)
	if s.activeID == id {
		s.activeID = ""
	}
	s.mu.Unlock()

	return s.repo.Delete(id)
}

// GetConnector gets or establishes the active Connector instance for a given connection ID.
func (s *Service) GetConnector(ctx context.Context, id string) (Connector, *ConnectionConfig, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if conn, exists := s.connectors[id]; exists {
		cfg := s.configs[id]
		return conn, &cfg, nil
	}

	cfg, err := s.repo.Get(id)
	if err != nil {
		return nil, nil, err
	}
	if cfg == nil {
		return nil, nil, ErrConnectionNotFound
	}

	conn, err := s.factory(cfg.Driver)
	if err != nil {
		return nil, nil, err
	}

	if err := conn.Connect(ctx, *cfg); err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to open connection %s: %w", id, err)
	}

	s.connectors[id] = conn
	s.configs[id] = *cfg
	return conn, cfg, nil
}

// SetActiveID sets the currently focused connection ID.
func (s *Service) SetActiveID(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeID = id
}

// GetActiveID returns the current focused connection ID.
func (s *Service) GetActiveID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.activeID
}

// CloseAll closes all cached active connectors.
func (s *Service) CloseAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, conn := range s.connectors {
		conn.Close()
		delete(s.connectors, id)
	}
}
