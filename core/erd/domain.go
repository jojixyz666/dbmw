package erd

import (
	"dbmw/core/explorer"
)

// Position represents 2D coordinates on canvas.
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// ERDNode represents a table in the ERD diagram.
type ERDNode struct {
	ID       string            `json:"id"`
	Label    string            `json:"label"`
	Schema   string            `json:"schema,omitempty"`
	Columns  []explorer.Column `json:"columns"`
	Position Position          `json:"position"`
}

// ERDEdge represents a relationship between two tables derived from foreign keys.
type ERDEdge struct {
	ID           string `json:"id"`
	SourceNode   string `json:"source"`
	SourceColumn string `json:"sourceHandle"`
	TargetNode   string `json:"target"`
	TargetColumn string `json:"targetHandle"`
	Label        string `json:"label,omitempty"`
	RelationType string `json:"relationType,omitempty"` // "one-to-many", etc.
}

// ERDGraph represents the complete ERD dataset ready for rendering.
type ERDGraph struct {
	Nodes []ERDNode `json:"nodes"`
	Edges []ERDEdge `json:"edges"`
}
