package erd

import (
	"context"
	"fmt"
	"math"
	"dbmw/core/explorer"
)

// ExplorerPort defines the explorer methods needed by ERDService.
type ExplorerPort interface {
	ListTables(ctx context.Context, schema string) ([]explorer.Table, error)
	ListColumns(ctx context.Context, schema, table string) ([]explorer.Column, error)
	ListForeignKeys(ctx context.Context, schema, table string) ([]explorer.ForeignKey, error)
}

// Service provides ERD graph generation.
type Service struct{}

// NewService instantiates ERDService.
func NewService() *Service {
	return &Service{}
}

// GenerateGraph builds table nodes and foreign key edges.
func (s *Service) GenerateGraph(ctx context.Context, explorerPort ExplorerPort, schema string) (*ERDGraph, error) {
	tables, err := explorerPort.ListTables(ctx, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables for schema %s: %w", schema, err)
	}

	nodes := make([]ERDNode, 0, len(tables))
	edges := make([]ERDEdge, 0)
	edgeSet := make(map[string]bool)

	// Layout parameters
	colsPerRow := int(math.Ceil(math.Sqrt(float64(len(tables)))))
	if colsPerRow < 1 {
		colsPerRow = 1
	}
	if colsPerRow > 4 {
		colsPerRow = 4
	}

	colSpacing := 320.0
	rowSpacing := 360.0

	for i, tbl := range tables {
		columns, err := explorerPort.ListColumns(ctx, schema, tbl.Name)
		if err != nil {
			columns = []explorer.Column{}
		}

		fks, err := explorerPort.ListForeignKeys(ctx, schema, tbl.Name)
		if err == nil {
			for _, fk := range fks {
				edgeKey := fmt.Sprintf("%s:%s->%s:%s", tbl.Name, fk.Column, fk.RefTable, fk.RefColumn)
				if !edgeSet[edgeKey] {
					edgeSet[edgeKey] = true
					edges = append(edges, ERDEdge{
						ID:           fmt.Sprintf("e_%s_%s", tbl.Name, fk.Name),
						SourceNode:   tbl.Name,
						SourceColumn: fk.Column,
						TargetNode:   fk.RefTable,
						TargetColumn: fk.RefColumn,
						Label:        fk.Name,
						RelationType: "many-to-one",
					})
				}
			}
		}

		gridX := float64(i % colsPerRow)
		gridY := float64(i / colsPerRow)

		nodes = append(nodes, ERDNode{
			ID:      tbl.Name,
			Label:   tbl.Name,
			Schema:  tbl.Schema,
			Columns: columns,
			Position: Position{
				X: 40.0 + gridX*colSpacing,
				Y: 40.0 + gridY*rowSpacing,
			},
		})
	}

	return &ERDGraph{
		Nodes: nodes,
		Edges: edges,
	}, nil
}
