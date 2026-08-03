package mcp

import (
	"context"
	"fmt"

	"dbmw/core/connection"
	"dbmw/core/explorer"
)

// SchemaTools provides read-only schema introspection tools.
type SchemaTools struct {
	explorerSvc *explorer.Service
}

func NewSchemaTools(expSvc *explorer.Service) *SchemaTools {
	return &SchemaTools{explorerSvc: expSvc}
}

func (t *SchemaTools) ReadSchemas(ctx context.Context, conn connection.Connector, database string) ([]explorer.Schema, error) {
	return t.explorerSvc.GetSchemas(ctx, conn, database)
}

func (t *SchemaTools) ReadTables(ctx context.Context, conn connection.Connector, schema string) ([]explorer.Table, error) {
	return t.explorerSvc.GetTables(ctx, conn, schema)
}

func (t *SchemaTools) ReadColumns(ctx context.Context, conn connection.Connector, schema, table string) ([]explorer.Column, error) {
	if table == "" {
		return nil, fmt.Errorf("table name is required")
	}
	return t.explorerSvc.GetColumns(ctx, conn, schema, table)
}
