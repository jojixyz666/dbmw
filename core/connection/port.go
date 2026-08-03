package connection

import (
	"context"
	"dbmw/core/data"
	"dbmw/core/explorer"
	"dbmw/core/query"
)

// Connector defines operations that any database driver adapter must fulfill.
type Connector interface {
	Connect(ctx context.Context, cfg ConnectionConfig) error
	Close() error
	Ping(ctx context.Context) error

	// Metadata & Exploration
	ListDatabases(ctx context.Context) ([]explorer.Database, error)
	ListSchemas(ctx context.Context, database string) ([]explorer.Schema, error)
	ListTables(ctx context.Context, schema string) ([]explorer.Table, error)
	ListColumns(ctx context.Context, schema, table string) ([]explorer.Column, error)
	ListIndexes(ctx context.Context, schema, table string) ([]explorer.Index, error)
	ListForeignKeys(ctx context.Context, schema, table string) ([]explorer.ForeignKey, error)
	ListViews(ctx context.Context, schema string) ([]explorer.View, error)

	// SQL Execution
	Execute(ctx context.Context, sqlQuery string) (*query.QueryResult, error)

	// Data operations
	BrowseRows(ctx context.Context, schema, table string, opts data.BrowseOptions) (*data.RowPage, error)
	InsertRow(ctx context.Context, schema, table string, values map[string]any) error
	UpdateRow(ctx context.Context, schema, table string, pk map[string]any, values map[string]any) error
	DeleteRow(ctx context.Context, schema, table string, pk map[string]any) error
}

// ConnectionRepository defines persistence methods for saved connections profile.
type ConnectionRepository interface {
	Get(id string) (*ConnectionConfig, error)
	GetAll() ([]ConnectionConfig, error)
	Save(cfg ConnectionConfig) error
	Delete(id string) error
}

// ConnectorFactory creates a new uninitialized Connector for the given driver.
type ConnectorFactory func(driver DriverType) (Connector, error)
