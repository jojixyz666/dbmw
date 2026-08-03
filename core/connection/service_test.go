package connection_test

import (
	"context"
	"testing"
	"dbmw/core/connection"
	"dbmw/core/data"
	"dbmw/core/explorer"
	"dbmw/core/query"
)

type fakeConnector struct {
	pingErr error
}

func (f *fakeConnector) Connect(ctx context.Context, cfg connection.ConnectionConfig) error {
	return f.pingErr
}
func (f *fakeConnector) Close() error { return nil }
func (f *fakeConnector) Ping(ctx context.Context) error {
	return f.pingErr
}
func (f *fakeConnector) ListDatabases(ctx context.Context) ([]explorer.Database, error) {
	return []explorer.Database{{Name: "testdb"}}, nil
}
func (f *fakeConnector) ListSchemas(ctx context.Context, database string) ([]explorer.Schema, error) {
	return []explorer.Schema{{Name: "public"}}, nil
}
func (f *fakeConnector) ListTables(ctx context.Context, schema string) ([]explorer.Table, error) {
	return []explorer.Table{{Name: "users", Type: "BASE TABLE"}}, nil
}
func (f *fakeConnector) ListColumns(ctx context.Context, schema, table string) ([]explorer.Column, error) {
	return []explorer.Column{{Name: "id", DataType: "integer", IsPrimaryKey: true}}, nil
}
func (f *fakeConnector) ListIndexes(ctx context.Context, schema, table string) ([]explorer.Index, error) {
	return nil, nil
}
func (f *fakeConnector) ListForeignKeys(ctx context.Context, schema, table string) ([]explorer.ForeignKey, error) {
	return nil, nil
}
func (f *fakeConnector) ListViews(ctx context.Context, schema string) ([]explorer.View, error) {
	return nil, nil
}
func (f *fakeConnector) Execute(ctx context.Context, sqlQuery string) (*query.QueryResult, error) {
	return &query.QueryResult{Columns: []string{"id"}, RowsAffected: 1}, nil
}
func (f *fakeConnector) BrowseRows(ctx context.Context, schema, table string, opts data.BrowseOptions) (*data.RowPage, error) {
	return &data.RowPage{TotalRows: 1}, nil
}
func (f *fakeConnector) InsertRow(ctx context.Context, schema, table string, values map[string]any) error {
	return nil
}
func (f *fakeConnector) UpdateRow(ctx context.Context, schema, table string, pk map[string]any, values map[string]any) error {
	return nil
}
func (f *fakeConnector) DeleteRow(ctx context.Context, schema, table string, pk map[string]any) error {
	return nil
}

type fakeConnectionRepo struct {
	items map[string]connection.ConnectionConfig
}

func newFakeRepo() *fakeConnectionRepo {
	return &fakeConnectionRepo{items: make(map[string]connection.ConnectionConfig)}
}

func (r *fakeConnectionRepo) Get(id string) (*connection.ConnectionConfig, error) {
	it, ok := r.items[id]
	if !ok {
		return nil, nil
	}
	return &it, nil
}

func (r *fakeConnectionRepo) GetAll() ([]connection.ConnectionConfig, error) {
	res := make([]connection.ConnectionConfig, 0, len(r.items))
	for _, it := range r.items {
		res = append(res, it)
	}
	return res, nil
}

func (r *fakeConnectionRepo) Save(cfg connection.ConnectionConfig) error {
	r.items[cfg.ID] = cfg
	return nil
}

func (r *fakeConnectionRepo) Delete(id string) error {
	delete(r.items, id)
	return nil
}

func TestConnectionService(t *testing.T) {
	repo := newFakeRepo()
	factory := func(driver connection.DriverType) (connection.Connector, error) {
		return &fakeConnector{}, nil
	}
	svc := connection.NewService(repo, factory)

	t.Run("Save and Get Connection", func(t *testing.T) {
		cfg := connection.ConnectionConfig{
			Name:   "Local Postgres",
			Driver: connection.DriverPostgres,
			Host:   "localhost",
			Port:   5432,
		}
		saved, err := svc.SaveConnection(cfg)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if saved.ID == "" {
			t.Fatal("expected generated ID")
		}

		got, err := svc.GetConnection(saved.ID)
		if err != nil {
			t.Fatalf("expected to get connection, got %v", err)
		}
		if got.Name != "Local Postgres" {
			t.Errorf("expected name 'Local Postgres', got %s", got.Name)
		}
	})

	t.Run("Test Connection Success", func(t *testing.T) {
		cfg := connection.ConnectionConfig{
			Name:   "SQLite File",
			Driver: connection.DriverSQLite,
			FilePath: ":memory:",
		}
		err := svc.TestConnection(context.Background(), cfg)
		if err != nil {
			t.Fatalf("expected test connection ok, got %v", err)
		}
	})

	t.Run("Test Connection Failure on Invalid Host", func(t *testing.T) {
		cfg := connection.ConnectionConfig{
			Name:   "Invalid Config",
			Driver: connection.DriverMySQL,
			Host:   "", // missing host
		}
		err := svc.TestConnection(context.Background(), cfg)
		if err == nil {
			t.Fatal("expected validation error on empty host")
		}
	})
}
