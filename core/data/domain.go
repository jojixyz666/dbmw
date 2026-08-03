package data

type FilterOperator string

const (
	OpEq         FilterOperator = "eq"
	OpNeq        FilterOperator = "neq"
	OpContains   FilterOperator = "contains"
	OpStartsWith FilterOperator = "starts_with"
	OpEndsWith   FilterOperator = "ends_with"
	OpGt         FilterOperator = "gt"
	OpGte        FilterOperator = "gte"
	OpLt         FilterOperator = "lt"
	OpLte        FilterOperator = "lte"
	OpIsNull     FilterOperator = "is_null"
	OpIsNotNull  FilterOperator = "is_not_null"
	OpIn         FilterOperator = "in"
)

// FilterClause specifies a condition on a column.
type FilterClause struct {
	Column   string         `json:"column"`
	Operator FilterOperator `json:"operator"`
	Value    any            `json:"value"`
}

// SortDirection specifies asc or desc.
type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

// SortClause specifies sorting order for a column.
type SortClause struct {
	Column    string        `json:"column"`
	Direction SortDirection `json:"direction"`
}

// BrowseOptions defines parameters for pagination, sorting, and filtering table rows.
type BrowseOptions struct {
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	Filters    []FilterClause `json:"filters,omitempty"`
	SortBy     *SortClause    `json:"sortBy,omitempty"`
	SearchTerm string         `json:"searchTerm,omitempty"`
}

// RowPage contains a page of data rows and metadata.
type RowPage struct {
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
	TotalRows  int64            `json:"totalRows"`
	TotalPages int              `json:"totalPages"`
	Columns    []string         `json:"columns"`
	Rows       []map[string]any `json:"rows"`
	PrimaryKey []string         `json:"primaryKey,omitempty"`
}
