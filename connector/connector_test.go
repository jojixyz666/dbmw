package connector_test

import (
	"context"
	"path/filepath"
	"testing"

	"dbmw/connector"
	"dbmw/core/connection"
	"dbmw/core/data"
)

func TestConnectorFactory(t *testing.T) {
	drivers := []connection.DriverType{
		connection.DriverPostgres,
		connection.DriverMySQL,
		connection.DriverMariaDB,
		connection.DriverSQLite,
	}

	for _, d := range drivers {
		t.Run(string(d), func(t *testing.T) {
			conn, err := connector.DefaultFactory(d)
			if err != nil {
				t.Fatalf("unexpected error for %s: %v", d, err)
			}
			if conn == nil {
				t.Fatalf("expected non-nil connector for %s", d)
			}
		})
	}
}

func TestSQLiteConnectorEndToEnd(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	conn, err := connector.DefaultFactory(connection.DriverSQLite)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()
	cfg := connection.ConnectionConfig{
		Name:     "SQLite Local",
		Driver:   connection.DriverSQLite,
		FilePath: dbPath,
	}

	if err := conn.Connect(ctx, cfg); err != nil {
		t.Fatalf("failed to connect sqlite: %v", err)
	}

	// 1. Create table with foreign key and index
	createDDL := `
	CREATE TABLE categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	CREATE TABLE articles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_id INTEGER,
		title TEXT NOT NULL,
		content TEXT,
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);
	CREATE INDEX idx_articles_title ON articles(title);
	CREATE VIEW v_articles AS SELECT id, title FROM articles;
	`
	if _, err := conn.Execute(ctx, createDDL); err != nil {
		t.Fatalf("failed to create tables: %v", err)
	}

	// 2. Test metadata listing
	tables, err := conn.ListTables(ctx, "main")
	if err != nil || len(tables) < 2 {
		t.Fatalf("expected at least 2 tables/views, got %v, err %v", tables, err)
	}

	cols, err := conn.ListColumns(ctx, "main", "articles")
	if err != nil || len(cols) != 4 {
		t.Fatalf("expected 4 columns in articles, got %v, err %v", cols, err)
	}

	idxs, err := conn.ListIndexes(ctx, "main", "articles")
	if err != nil || len(idxs) == 0 {
		t.Fatalf("expected indexes on articles, got %v, err %v", idxs, err)
	}

	fks, err := conn.ListForeignKeys(ctx, "main", "articles")
	if err != nil || len(fks) != 1 {
		t.Fatalf("expected 1 foreign key from articles to categories, got %v, err %v", fks, err)
	}

	views, err := conn.ListViews(ctx, "main")
	if err != nil || len(views) == 0 {
		t.Fatalf("expected view v_articles, got %v, err %v", views, err)
	}

	// 3. Test insert, update, browse, delete
	if err := conn.InsertRow(ctx, "main", "categories", map[string]any{"name": "Tech"}); err != nil {
		t.Fatalf("failed to insert category: %v", err)
	}

	if err := conn.InsertRow(ctx, "main", "articles", map[string]any{
		"category_id": 1,
		"title":       "Testing Golang & SQLite",
		"content":     "Full stack local database manager",
	}); err != nil {
		t.Fatalf("failed to insert article: %v", err)
	}

	browsePage, err := conn.BrowseRows(ctx, "main", "articles", data.BrowseOptions{
		Page:     1,
		PageSize: 10,
	})
	if err != nil || browsePage.TotalRows != 1 {
		t.Fatalf("expected 1 row in browse, got %+v, err %v", browsePage, err)
	}

	// Update row
	if err := conn.UpdateRow(ctx, "main", "articles", map[string]any{"id": 1}, map[string]any{"title": "Updated Title"}); err != nil {
		t.Fatalf("failed to update article: %v", err)
	}

	// Query result
	qRes, err := conn.Execute(ctx, "SELECT title FROM articles WHERE id = 1;")
	if err != nil || len(qRes.Rows) == 0 || qRes.Rows[0]["title"] != "Updated Title" {
		t.Fatalf("failed to query updated record: %+v, err %v", qRes, err)
	}

	// Delete row
	if err := conn.DeleteRow(ctx, "main", "articles", map[string]any{"id": 1}); err != nil {
		t.Fatalf("failed to delete article: %v", err)
	}

	browseAfter, _ := conn.BrowseRows(ctx, "main", "articles", data.BrowseOptions{Page: 1, PageSize: 10})
	if browseAfter.TotalRows != 0 {
		t.Errorf("expected 0 rows after delete, got %d", browseAfter.TotalRows)
	}
}
