package explorer

// Database represents a physical or logical database server catalog.
type Database struct {
	Name string `json:"name"`
}

// Schema represents a database namespace or schema.
type Schema struct {
	Name string `json:"name"`
}

// Table represents a database table metadata.
type Table struct {
	Name      string `json:"name"`
	Schema    string `json:"schema,omitempty"`
	Type      string `json:"type"` // "BASE TABLE", "VIEW", "SYSTEM TABLE"
	RowCount  int64  `json:"rowCount,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

// Column represents a column metadata in a table.
type Column struct {
	Name         string  `json:"name"`
	DataType     string  `json:"dataType"`
	IsNullable   bool    `json:"isNullable"`
	IsPrimaryKey bool    `json:"isPrimaryKey"`
	DefaultValue *string `json:"defaultValue,omitempty"`
	Comment      string  `json:"comment,omitempty"`
	Extra        string  `json:"extra,omitempty"`
	Position     int     `json:"position"`
}

// Index represents an index on a table.
type Index struct {
	Name      string   `json:"name"`
	Columns   []string `json:"columns"`
	IsUnique  bool     `json:"isUnique"`
	IsPrimary bool     `json:"isPrimary"`
	Type      string   `json:"type,omitempty"` // BTREE, HASH, GIN, etc.
}

// ForeignKey represents a foreign key constraint between tables.
type ForeignKey struct {
	Name         string `json:"name"`
	Column       string `json:"column"`
	RefTable     string `json:"refTable"`
	RefColumn    string `json:"refColumn"`
	RefSchema    string `json:"refSchema,omitempty"`
	OnDelete     string `json:"onDelete,omitempty"`
	OnUpdate     string `json:"onUpdate,omitempty"`
}

// View represents a database view.
type View struct {
	Name       string `json:"name"`
	Schema     string `json:"schema,omitempty"`
	Definition string `json:"definition,omitempty"`
}
