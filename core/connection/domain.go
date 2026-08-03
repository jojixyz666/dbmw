package connection

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// DriverType represents supported database engines.
type DriverType string

const (
	DriverPostgres DriverType = "postgres"
	DriverMySQL    DriverType = "mysql"
	DriverMariaDB  DriverType = "mariadb"
	DriverSQLite   DriverType = "sqlite"
)

var (
	ErrConnectionNotFound = errors.New("connection not found")
	ErrInvalidConfig      = errors.New("invalid connection configuration")
	ErrUnsupportedDriver  = errors.New("unsupported database driver")
)

// ConnectionConfig defines the connection parameters for a database.
type ConnectionConfig struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Driver    DriverType        `json:"driver"`
	Host      string            `json:"host,omitempty"`
	Port      int               `json:"port,omitempty"`
	User      string            `json:"user,omitempty"`
	Password  string            `json:"password,omitempty"`
	Database  string            `json:"database,omitempty"`
	FilePath  string            `json:"filePath,omitempty"` // For SQLite
	SSLMode   string            `json:"sslMode,omitempty"`
	Params    map[string]string `json:"params,omitempty"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

// Validate checks if the connection configuration has required fields.
func (c *ConnectionConfig) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidConfig)
	}
	switch c.Driver {
	case DriverPostgres, DriverMySQL, DriverMariaDB:
		if strings.TrimSpace(c.Host) == "" {
			return fmt.Errorf("%w: host is required for driver %s", ErrInvalidConfig, c.Driver)
		}
		if c.Port <= 0 {
			switch c.Driver {
			case DriverPostgres:
				c.Port = 5432
			case DriverMySQL, DriverMariaDB:
				c.Port = 3306
			}
		}
	case DriverSQLite:
		if strings.TrimSpace(c.FilePath) == "" && strings.TrimSpace(c.Database) == "" {
			return fmt.Errorf("%w: filePath or database is required for sqlite", ErrInvalidConfig)
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedDriver, c.Driver)
	}
	return nil
}
