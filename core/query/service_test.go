package query_test

import (
	"context"
	"strings"
	"testing"
	"dbmw/core/query"
)

type fakeQueryConnector struct {
	fail bool
}

func (f *fakeQueryConnector) Execute(ctx context.Context, sqlQuery string) (*query.QueryResult, error) {
	if f.fail {
		return nil, context.DeadlineExceeded
	}
	return &query.QueryResult{
		Columns:      []string{"id", "title"},
		ColumnTypes:  []string{"int", "text"},
		Rows:         []map[string]any{{"id": 1, "title": "Hello"}},
		RowsAffected: 1,
	}, nil
}

type fakeHistoryRepo struct {
	items []query.QueryHistoryItem
}

func (f *fakeHistoryRepo) Save(ctx context.Context, item query.QueryHistoryItem) error {
	f.items = append(f.items, item)
	return nil
}

func (f *fakeHistoryRepo) List(ctx context.Context, connID string, limit int) ([]query.QueryHistoryItem, error) {
	return f.items, nil
}

func (f *fakeHistoryRepo) Clear(ctx context.Context, connID string) error {
	f.items = nil
	return nil
}

func TestQueryService(t *testing.T) {
	repo := &fakeHistoryRepo{}
	svc := query.NewService(repo)

	t.Run("Execute Success and History Logged", func(t *testing.T) {
		conn := &fakeQueryConnector{}
		res, err := svc.Execute(context.Background(), conn, "conn_1", "Local", "SELECT * FROM posts")
		if err != nil {
			t.Fatalf("expected success, got %v", err)
		}
		if len(res.Rows) != 1 {
			t.Fatalf("expected 1 row, got %d", len(res.Rows))
		}

		hist, err := svc.GetHistory(context.Background(), "conn_1", 10)
		if err != nil || len(hist) != 1 {
			t.Fatalf("expected 1 history item, got %v", err)
		}
		if hist[0].Status != "success" {
			t.Errorf("expected success status, got %s", hist[0].Status)
		}
	})

	t.Run("Export CSV", func(t *testing.T) {
		res := &query.QueryResult{
			Columns: []string{"id", "username"},
			Rows: []map[string]any{
				{"id": 10, "username": "alice"},
				{"id": 20, "username": "bob"},
			},
		}
		csvBytes, err := svc.ExportCSV(res)
		if err != nil {
			t.Fatalf("failed to export csv: %v", err)
		}
		csvStr := string(csvBytes)
		if !strings.Contains(csvStr, "id,username") || !strings.Contains(csvStr, "alice") {
			t.Errorf("unexpected csv output: %s", csvStr)
		}
	})
}
