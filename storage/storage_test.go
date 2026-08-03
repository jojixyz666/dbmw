package storage_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"
	"dbmw/core/connection"
	"dbmw/core/query"
	"dbmw/storage"
)

func TestStorage(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("ConfigStore save and load", func(t *testing.T) {
		p := filepath.Join(tempDir, "config.json")
		store, err := storage.NewConfigStore(p)
		if err != nil {
			t.Fatal(err)
		}

		cfg := storage.AppConfig{
			Theme:           "light",
			DefaultPageSize: 50,
			ServerPort:      9090,
		}
		if err := store.Save(cfg); err != nil {
			t.Fatalf("failed to save config: %v", err)
		}

		loaded, err := store.Get()
		if err != nil {
			t.Fatalf("failed to get config: %v", err)
		}
		if loaded.Theme != "light" || loaded.DefaultPageSize != 50 || loaded.ServerPort != 9090 {
			t.Errorf("mismatch config values: %+v", loaded)
		}
	})

	t.Run("ConnectionStore save, get, list, delete", func(t *testing.T) {
		p := filepath.Join(tempDir, "connections.json")
		store, err := storage.NewConnectionStore(p)
		if err != nil {
			t.Fatal(err)
		}

		connItem := connection.ConnectionConfig{
			ID:     "conn_sqlite_test",
			Name:   "SQLite File",
			Driver: connection.DriverSQLite,
			FilePath: ":memory:",
		}

		if err := store.Save(connItem); err != nil {
			t.Fatalf("failed to save connection: %v", err)
		}

		all, err := store.GetAll()
		if err != nil || len(all) != 1 {
			t.Fatalf("expected 1 item, got %d, err %v", len(all), err)
		}

		got, err := store.Get("conn_sqlite_test")
		if err != nil || got == nil || got.Name != "SQLite File" {
			t.Fatalf("failed to get single connection: %v", got)
		}

		if err := store.Delete("conn_sqlite_test"); err != nil {
			t.Fatalf("failed to delete connection: %v", err)
		}

		allAfter, _ := store.GetAll()
		if len(allAfter) != 0 {
			t.Errorf("expected 0 items after delete, got %d", len(allAfter))
		}
	})

	t.Run("HistoryStore log, list, clear", func(t *testing.T) {
		p := filepath.Join(tempDir, "history.db")
		store, err := storage.NewHistoryStore(p)
		if err != nil {
			t.Fatal(err)
		}
		defer store.Close()

		ctx := context.Background()
		item := query.QueryHistoryItem{
			ID:           "hist_1",
			ConnectionID: "conn_123",
			Query:        "SELECT 1;",
			ExecutedAt:   time.Now(),
			DurationMs:   12.5,
			Status:       "success",
			RowsAffected: 1,
		}

		if err := store.Save(ctx, item); err != nil {
			t.Fatalf("failed to save history item: %v", err)
		}

		list, err := store.List(ctx, "conn_123", 10)
		if err != nil || len(list) != 1 {
			t.Fatalf("expected 1 history item, got %d, err %v", len(list), err)
		}

		if list[0].Query != "SELECT 1;" {
			t.Errorf("unexpected query: %s", list[0].Query)
		}

		if err := store.Clear(ctx, "conn_123"); err != nil {
			t.Fatalf("failed to clear history: %v", err)
		}

		listAfter, _ := store.List(ctx, "conn_123", 10)
		if len(listAfter) != 0 {
			t.Errorf("expected 0 history items after clear, got %d", len(listAfter))
		}
	})
}
