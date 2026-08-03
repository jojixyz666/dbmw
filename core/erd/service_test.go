package erd_test

import (
	"context"
	"testing"
	"dbmw/core/erd"
	"dbmw/core/explorer"
)

type fakeErdExplorer struct{}

func (f *fakeErdExplorer) ListTables(ctx context.Context, schema string) ([]explorer.Table, error) {
	return []explorer.Table{
		{Name: "users", Schema: "public"},
		{Name: "posts", Schema: "public"},
	}, nil
}

func (f *fakeErdExplorer) ListColumns(ctx context.Context, schema, table string) ([]explorer.Column, error) {
	if table == "users" {
		return []explorer.Column{
			{Name: "id", DataType: "bigint", IsPrimaryKey: true},
			{Name: "email", DataType: "varchar"},
		}, nil
	}
	return []explorer.Column{
		{Name: "id", DataType: "bigint", IsPrimaryKey: true},
		{Name: "user_id", DataType: "bigint"},
		{Name: "title", DataType: "varchar"},
	}, nil
}

func (f *fakeErdExplorer) ListForeignKeys(ctx context.Context, schema, table string) ([]explorer.ForeignKey, error) {
	if table == "posts" {
		return []explorer.ForeignKey{
			{
				Name:      "fk_posts_users",
				Column:    "user_id",
				RefTable:  "users",
				RefColumn: "id",
			},
		}, nil
	}
	return nil, nil
}

func TestERDService(t *testing.T) {
	svc := erd.NewService()
	fake := &fakeErdExplorer{}

	graph, err := svc.GenerateGraph(context.Background(), fake, "public")
	if err != nil {
		t.Fatalf("expected graph generation to succeed: %v", err)
	}

	if len(graph.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(graph.Nodes))
	}
	if len(graph.Edges) != 1 {
		t.Fatalf("expected 1 fk relationship edge, got %d", len(graph.Edges))
	}
	if graph.Edges[0].SourceNode != "posts" || graph.Edges[0].TargetNode != "users" {
		t.Errorf("unexpected edge relation: %+v", graph.Edges[0])
	}
}
