package data_test

import (
	"context"
	"testing"
	"dbmw/core/data"
)

type fakeDataConnector struct{}

func (f *fakeDataConnector) BrowseRows(ctx context.Context, schema, table string, opts data.BrowseOptions) (*data.RowPage, error) {
	return &data.RowPage{
		Page:       opts.Page,
		PageSize:   opts.PageSize,
		TotalRows:  100,
		TotalPages: 10,
		Columns:    []string{"id", "title"},
		Rows:       []map[string]any{{"id": 1, "title": "My Post"}},
		PrimaryKey: []string{"id"},
	}, nil
}

func (f *fakeDataConnector) InsertRow(ctx context.Context, schema, table string, values map[string]any) error {
	return nil
}

func (f *fakeDataConnector) UpdateRow(ctx context.Context, schema, table string, pk map[string]any, values map[string]any) error {
	return nil
}

func (f *fakeDataConnector) DeleteRow(ctx context.Context, schema, table string, pk map[string]any) error {
	return nil
}

func TestDataService(t *testing.T) {
	svc := data.NewService()
	conn := &fakeDataConnector{}
	ctx := context.Background()

	t.Run("Browse with Default Options", func(t *testing.T) {
		page, err := svc.Browse(ctx, conn, "public", "posts", data.BrowseOptions{})
		if err != nil {
			t.Fatalf("expected no error: %v", err)
		}
		if page.Page != 1 || page.PageSize != 25 {
			t.Errorf("expected page 1 and pageSize 25, got %d and %d", page.Page, page.PageSize)
		}
	})

	t.Run("Insert Valid Values", func(t *testing.T) {
		err := svc.Insert(ctx, conn, "public", "posts", map[string]any{"title": "Awesome"})
		if err != nil {
			t.Fatalf("expected successful insert: %v", err)
		}
	})

	t.Run("Update without Primary Key Fails", func(t *testing.T) {
		err := svc.Update(ctx, conn, "public", "posts", nil, map[string]any{"title": "New Title"})
		if err == nil {
			t.Fatal("expected error without primary key")
		}
	})
}
