package query

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"
)

// Service provides query execution and history management.
type Service struct {
	historyRepo HistoryRepository
}

// NewService instantiates QueryService.
func NewService(repo HistoryRepository) *Service {
	return &Service{historyRepo: repo}
}

// Execute runs SQL query on the given connector and logs it to history.
func (s *Service) Execute(ctx context.Context, conn Connector, connID, connName, sqlQuery string) (*QueryResult, error) {
	start := time.Now()
	res, err := conn.Execute(ctx, sqlQuery)
	duration := time.Since(start).Seconds() * 1000.0

	historyItem := QueryHistoryItem{
		ID:             fmt.Sprintf("hist_%d", time.Now().UnixNano()),
		ConnectionID:   connID,
		ConnectionName: connName,
		Query:          sqlQuery,
		ExecutedAt:     start,
		DurationMs:     duration,
	}

	if err != nil {
		historyItem.Status = "error"
		historyItem.ErrorMessage = err.Error()
		if s.historyRepo != nil {
			_ = s.historyRepo.Save(ctx, historyItem)
		}
		return &QueryResult{
			ExecutionTimeMs: duration,
			Error:           err.Error(),
		}, err
	}

	if res == nil {
		res = &QueryResult{}
	}
	res.ExecutionTimeMs = duration

	historyItem.Status = "success"
	historyItem.RowsAffected = res.RowsAffected
	if s.historyRepo != nil {
		_ = s.historyRepo.Save(ctx, historyItem)
	}

	return res, nil
}

// GetHistory returns recent historical query logs.
func (s *Service) GetHistory(ctx context.Context, connectionID string, limit int) ([]QueryHistoryItem, error) {
	if s.historyRepo == nil {
		return []QueryHistoryItem{}, nil
	}
	if limit <= 0 {
		limit = 50
	}
	return s.historyRepo.List(ctx, connectionID, limit)
}

// ClearHistory deletes history for a connection.
func (s *Service) ClearHistory(ctx context.Context, connectionID string) error {
	if s.historyRepo == nil {
		return nil
	}
	return s.historyRepo.Clear(ctx, connectionID)
}

// ExportCSV formats rows from a query result into CSV format.
func (s *Service) ExportCSV(result *QueryResult) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Header
	if err := writer.Write(result.Columns); err != nil {
		return nil, err
	}

	// Rows
	for _, row := range result.Rows {
		record := make([]string, len(result.Columns))
		for i, col := range result.Columns {
			val := row[col]
			if val == nil {
				record[i] = ""
			} else {
				record[i] = fmt.Sprintf("%v", val)
			}
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ExportJSON formats rows from a query result into clean JSON.
func (s *Service) ExportJSON(result *QueryResult) ([]byte, error) {
	return json.MarshalIndent(result.Rows, "", "  ")
}
