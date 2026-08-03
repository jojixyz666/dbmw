package explorer

import (
	"context"
)

// Connector is the consumer port needed by ExplorerService.
type Connector interface {
	ListDatabases(ctx context.Context) ([]Database, error)
	ListSchemas(ctx context.Context, database string) ([]Schema, error)
	ListTables(ctx context.Context, schema string) ([]Table, error)
	ListColumns(ctx context.Context, schema, table string) ([]Column, error)
	ListIndexes(ctx context.Context, schema, table string) ([]Index, error)
	ListForeignKeys(ctx context.Context, schema, table string) ([]ForeignKey, error)
	ListViews(ctx context.Context, schema string) ([]View, error)
}

// Service provides database metadata exploration capabilities.
type Service struct{}

// NewService instantiates ExplorerService.
func NewService() *Service {
	return &Service{}
}

// GetDatabases returns catalog/databases on current connector.
func (s *Service) GetDatabases(ctx context.Context, conn Connector) ([]Database, error) {
	return conn.ListDatabases(ctx)
}

// GetSchemas returns schemas for a database.
func (s *Service) GetSchemas(ctx context.Context, conn Connector, database string) ([]Schema, error) {
	return conn.ListSchemas(ctx, database)
}

// GetTables returns table metadata in a schema.
func (s *Service) GetTables(ctx context.Context, conn Connector, schema string) ([]Table, error) {
	return conn.ListTables(ctx, schema)
}

// GetTableDetails returns columns, indexes, foreign keys, and metadata for a single table.
func (s *Service) GetTableDetails(ctx context.Context, conn Connector, schema, table string) (map[string]any, error) {
	cols, err := conn.ListColumns(ctx, schema, table)
	if err != nil {
		return nil, err
	}
	idxs, err := conn.ListIndexes(ctx, schema, table)
	if err != nil {
		return nil, err
	}
	fks, err := conn.ListForeignKeys(ctx, schema, table)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"table":       table,
		"schema":      schema,
		"columns":     cols,
		"indexes":     idxs,
		"foreignKeys": fks,
	}, nil
}

// GetColumns returns columns for a table.
func (s *Service) GetColumns(ctx context.Context, conn Connector, schema, table string) ([]Column, error) {
	return conn.ListColumns(ctx, schema, table)
}

// GetIndexes returns indexes for a table.
func (s *Service) GetIndexes(ctx context.Context, conn Connector, schema, table string) ([]Index, error) {
	return conn.ListIndexes(ctx, schema, table)
}

// GetForeignKeys returns foreign keys for a table.
func (s *Service) GetForeignKeys(ctx context.Context, conn Connector, schema, table string) ([]ForeignKey, error) {
	return conn.ListForeignKeys(ctx, schema, table)
}

// GetViews returns database views.
func (s *Service) GetViews(ctx context.Context, conn Connector, schema string) ([]View, error) {
	return conn.ListViews(ctx, schema)
}
