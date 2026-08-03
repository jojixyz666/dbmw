package explorer_test

import (
	"context"
	"testing"
	"dbmw/core/explorer"
)

type fakeExplorerConnector struct{}

func (f *fakeExplorerConnector) ListDatabases(ctx context.Context) ([]explorer.Database, error) {
	return []explorer.Database{{Name: "master_db"}}, nil
}
func (f *fakeExplorerConnector) ListSchemas(ctx context.Context, database string) ([]explorer.Schema, error) {
	return []explorer.Schema{{Name: "public"}}, nil
}
func (f *fakeExplorerConnector) ListTables(ctx context.Context, schema string) ([]explorer.Table, error) {
	return []explorer.Table{
		{Name: "users", Schema: "public", Type: "BASE TABLE"},
		{Name: "posts", Schema: "public", Type: "BASE TABLE"},
	}, nil
}
func (f *fakeExplorerConnector) ListColumns(ctx context.Context, schema, table string) ([]explorer.Column, error) {
	return []explorer.Column{
		{Name: "id", DataType: "bigint", IsPrimaryKey: true, Position: 1},
		{Name: "name", DataType: "varchar(255)", IsNullable: false, Position: 2},
	}, nil
}
func (f *fakeExplorerConnector) ListIndexes(ctx context.Context, schema, table string) ([]explorer.Index, error) {
	return []explorer.Index{
		{Name: "users_pkey", Columns: []string{"id"}, IsPrimary: true, IsUnique: true},
	}, nil
}
func (f *fakeExplorerConnector) ListForeignKeys(ctx context.Context, schema, table string) ([]explorer.ForeignKey, error) {
	return []explorer.ForeignKey{}, nil
}
func (f *fakeExplorerConnector) ListViews(ctx context.Context, schema string) ([]explorer.View, error) {
	return []explorer.View{}, nil
}

func TestExplorerService(t *testing.T) {
	svc := explorer.NewService()
	conn := &fakeExplorerConnector{}
	ctx := context.Background()

	t.Run("GetDatabases", func(t *testing.T) {
		dbs, err := svc.GetDatabases(ctx, conn)
		if err != nil || len(dbs) != 1 {
			t.Fatalf("unexpected databases: %v, %v", dbs, err)
		}
	})

	t.Run("GetTables", func(t *testing.T) {
		tables, err := svc.GetTables(ctx, conn, "public")
		if err != nil || len(tables) != 2 {
			t.Fatalf("unexpected tables: %v, %v", tables, err)
		}
	})

	t.Run("GetTableDetails", func(t *testing.T) {
		details, err := svc.GetTableDetails(ctx, conn, "public", "users")
		if err != nil {
			t.Fatalf("expected table details, got err %v", err)
		}
		if details["table"] != "users" {
			t.Errorf("expected table name users, got %v", details["table"])
		}
	})
}
