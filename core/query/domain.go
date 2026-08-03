package query

import (
	"time"
)

// CellValue represents a formatted or raw SQL query result cell value.
type QueryResult struct {
	Columns         []string         `json:"columns"`
	ColumnTypes     []string         `json:"columnTypes"`
	Rows            []map[string]any `json:"rows"`
	RowsAffected    int64            `json:"rowsAffected"`
	ExecutionTimeMs float64          `json:"executionTimeMs"`
	Error           string           `json:"error,omitempty"`
}

// QueryHistoryItem represents a saved historical query run.
type QueryHistoryItem struct {
	ID              string    `json:"id"`
	ConnectionID    string    `json:"connectionId"`
	ConnectionName  string    `json:"connectionName,omitempty"`
	Query           string    `json:"query"`
	ExecutedAt      time.Time `json:"executedAt"`
	DurationMs      float64   `json:"durationMs"`
	Status          string    `json:"status"` // "success", "error"
	RowsAffected    int64     `json:"rowsAffected"`
	ErrorMessage    string    `json:"errorMessage,omitempty"`
}
