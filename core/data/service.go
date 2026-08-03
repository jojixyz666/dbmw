package data

import (
	"context"
	"errors"
)

var (
	ErrEmptyPrimaryKey = errors.New("primary key is required for updating or deleting a row")
	ErrEmptyValues     = errors.New("no values provided for row insert or update")
)

// Service provides row browsing and CRUD manipulation operations.
type Service struct{}

// NewService instantiates DataService.
func NewService() *Service {
	return &Service{}
}

// Browse retrieves a paginated set of rows with filter and sorting capabilities.
func (s *Service) Browse(ctx context.Context, conn Connector, schema, table string, opts BrowseOptions) (*RowPage, error) {
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 25
	}
	if opts.PageSize > 1000 {
		opts.PageSize = 1000
	}
	return conn.BrowseRows(ctx, schema, table, opts)
}

// Insert adds a new row to table.
func (s *Service) Insert(ctx context.Context, conn Connector, schema, table string, values map[string]any) error {
	if len(values) == 0 {
		return ErrEmptyValues
	}
	return conn.InsertRow(ctx, schema, table, values)
}

// Update modifies an existing row identified by its primary key.
func (s *Service) Update(ctx context.Context, conn Connector, schema, table string, pk map[string]any, values map[string]any) error {
	if len(pk) == 0 {
		return ErrEmptyPrimaryKey
	}
	if len(values) == 0 {
		return ErrEmptyValues
	}
	return conn.UpdateRow(ctx, schema, table, pk, values)
}

// Delete removes a row identified by primary key.
func (s *Service) Delete(ctx context.Context, conn Connector, schema, table string, pk map[string]any) error {
	if len(pk) == 0 {
		return ErrEmptyPrimaryKey
	}
	return conn.DeleteRow(ctx, schema, table, pk)
}
