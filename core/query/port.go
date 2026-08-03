package query

import (
	"context"
)

// Connector is the database execution port needed by QueryService.
type Connector interface {
	Execute(ctx context.Context, sqlQuery string) (*QueryResult, error)
}

// HistoryRepository defines persistence operations for historical query execution logs.
type HistoryRepository interface {
	Save(ctx context.Context, item QueryHistoryItem) error
	List(ctx context.Context, connectionID string, limit int) ([]QueryHistoryItem, error)
	Clear(ctx context.Context, connectionID string) error
}
