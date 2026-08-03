package mariadb

import (
	"context"
	"dbmw/connector/mysql"
	"dbmw/core/connection"
	"dbmw/core/data"
	"dbmw/core/explorer"
	"dbmw/core/query"
)

// Connector implements core/connection.Connector for MariaDB, wrapping MySQL protocol compatibility.
type Connector struct {
	inner connection.Connector
}

// NewConnector creates an uninitialized MariaDB connector.
func NewConnector() connection.Connector {
	return &Connector{
		inner: mysql.NewConnector(),
	}
}

func (c *Connector) Connect(ctx context.Context, cfg connection.ConnectionConfig) error {
	return c.inner.Connect(ctx, cfg)
}

func (c *Connector) Close() error {
	return c.inner.Close()
}

func (c *Connector) Ping(ctx context.Context) error {
	return c.inner.Ping(ctx)
}

func (c *Connector) ListDatabases(ctx context.Context) ([]explorer.Database, error) {
	return c.inner.ListDatabases(ctx)
}

func (c *Connector) ListSchemas(ctx context.Context, database string) ([]explorer.Schema, error) {
	return c.inner.ListSchemas(ctx, database)
}

func (c *Connector) ListTables(ctx context.Context, schema string) ([]explorer.Table, error) {
	return c.inner.ListTables(ctx, schema)
}

func (c *Connector) ListColumns(ctx context.Context, schema, table string) ([]explorer.Column, error) {
	return c.inner.ListColumns(ctx, schema, table)
}

func (c *Connector) ListIndexes(ctx context.Context, schema, table string) ([]explorer.Index, error) {
	return c.inner.ListIndexes(ctx, schema, table)
}

func (c *Connector) ListForeignKeys(ctx context.Context, schema, table string) ([]explorer.ForeignKey, error) {
	return c.inner.ListForeignKeys(ctx, schema, table)
}

func (c *Connector) ListViews(ctx context.Context, schema string) ([]explorer.View, error) {
	return c.inner.ListViews(ctx, schema)
}

func (c *Connector) Execute(ctx context.Context, sqlQuery string) (*query.QueryResult, error) {
	return c.inner.Execute(ctx, sqlQuery)
}

func (c *Connector) BrowseRows(ctx context.Context, schema, table string, opts data.BrowseOptions) (*data.RowPage, error) {
	return c.inner.BrowseRows(ctx, schema, table, opts)
}

func (c *Connector) InsertRow(ctx context.Context, schema, table string, values map[string]any) error {
	return c.inner.InsertRow(ctx, schema, table, values)
}

func (c *Connector) UpdateRow(ctx context.Context, schema, table string, pk map[string]any, values map[string]any) error {
	return c.inner.UpdateRow(ctx, schema, table, pk, values)
}

func (c *Connector) DeleteRow(ctx context.Context, schema, table string, pk map[string]any) error {
	return c.inner.DeleteRow(ctx, schema, table, pk)
}
