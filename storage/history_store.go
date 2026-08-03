package storage

import (
	"context"
	"database/sql"
	"path/filepath"
	"time"

	"dbmw/core/query"
	_ "modernc.org/sqlite"
)

// HistoryStore implements query.HistoryRepository backed by SQLite (~/.dbmw/history.db).
type HistoryStore struct {
	db *sql.DB
}

// NewHistoryStore initializes connection to the SQLite history file.
func NewHistoryStore(customPath ...string) (*HistoryStore, error) {
	var targetPath string
	if len(customPath) > 0 && customPath[0] != "" {
		targetPath = customPath[0]
	} else {
		dir, err := GetDbmwDir()
		if err != nil {
			return nil, err
		}
		targetPath = filepath.Join(dir, "history.db")
	}

	db, err := sql.Open("sqlite", targetPath)
	if err != nil {
		return nil, err
	}

	// Run migration
	schema := `
	CREATE TABLE IF NOT EXISTS query_history (
		id TEXT PRIMARY KEY,
		connection_id TEXT NOT NULL,
		connection_name TEXT,
		query TEXT NOT NULL,
		executed_at DATETIME NOT NULL,
		duration_ms REAL NOT NULL,
		status TEXT NOT NULL,
		rows_affected INTEGER DEFAULT 0,
		error_message TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_query_history_conn ON query_history(connection_id, executed_at DESC);
	`
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}

	return &HistoryStore{db: db}, nil
}

// Close releases the DB connection pool.
func (h *HistoryStore) Close() error {
	if h.db != nil {
		return h.db.Close()
	}
	return nil
}

// Save records a query execution log.
func (h *HistoryStore) Save(ctx context.Context, item query.QueryHistoryItem) error {
	stmt := `
	INSERT INTO query_history (id, connection_id, connection_name, query, executed_at, duration_ms, status, rows_affected, error_message)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := h.db.ExecContext(ctx, stmt,
		item.ID,
		item.ConnectionID,
		item.ConnectionName,
		item.Query,
		item.ExecutedAt.Format(time.RFC3339Nano),
		item.DurationMs,
		item.Status,
		item.RowsAffected,
		item.ErrorMessage,
	)
	return err
}

// List returns recent query history for a connection or all connections if connectionID is "".
func (h *HistoryStore) List(ctx context.Context, connectionID string, limit int) ([]query.QueryHistoryItem, error) {
	if limit <= 0 {
		limit = 50
	}

	var rows *sql.Rows
	var err error

	if connectionID != "" {
		stmt := `
		SELECT id, connection_id, connection_name, query, executed_at, duration_ms, status, rows_affected, COALESCE(error_message, '')
		FROM query_history
		WHERE connection_id = ?
		ORDER BY executed_at DESC
		LIMIT ?
		`
		rows, err = h.db.QueryContext(ctx, stmt, connectionID, limit)
	} else {
		stmt := `
		SELECT id, connection_id, connection_name, query, executed_at, duration_ms, status, rows_affected, COALESCE(error_message, '')
		FROM query_history
		ORDER BY executed_at DESC
		LIMIT ?
		`
		rows, err = h.db.QueryContext(ctx, stmt, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []query.QueryHistoryItem
	for rows.Next() {
		var item query.QueryHistoryItem
		var executedAtStr string
		var connName sql.NullString
		if err := rows.Scan(
			&item.ID,
			&item.ConnectionID,
			&connName,
			&item.Query,
			&executedAtStr,
			&item.DurationMs,
			&item.Status,
			&item.RowsAffected,
			&item.ErrorMessage,
		); err != nil {
			return nil, err
		}

		if connName.Valid {
			item.ConnectionName = connName.String
		}
		if t, err := time.Parse(time.RFC3339Nano, executedAtStr); err == nil {
			item.ExecutedAt = t
		} else if t, err := time.Parse("2006-01-02 15:04:05", executedAtStr); err == nil {
			item.ExecutedAt = t
		}

		result = append(result, item)
	}

	return result, rows.Err()
}

// Clear removes history logs for a specific connection or all.
func (h *HistoryStore) Clear(ctx context.Context, connectionID string) error {
	if connectionID != "" {
		_, err := h.db.ExecContext(ctx, "DELETE FROM query_history WHERE connection_id = ?", connectionID)
		return err
	}
	_, err := h.db.ExecContext(ctx, "DELETE FROM query_history")
	return err
}
