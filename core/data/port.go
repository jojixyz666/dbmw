package data

import (
	"context"
)

// Connector is the data manipulation port required by DataService.
type Connector interface {
	BrowseRows(ctx context.Context, schema, table string, opts BrowseOptions) (*RowPage, error)
	InsertRow(ctx context.Context, schema, table string, values map[string]any) error
	UpdateRow(ctx context.Context, schema, table string, pk map[string]any, values map[string]any) error
	DeleteRow(ctx context.Context, schema, table string, pk map[string]any) error
}
