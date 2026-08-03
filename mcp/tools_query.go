package mcp

import (
	"context"
	"fmt"

	"dbmw/core/connection"
	"dbmw/core/query"
)

// QueryTools provides read-only query execution tools guarded by Guard.
type QueryTools struct {
	querySvc *query.Service
}

func NewQueryTools(qSvc *query.Service) *QueryTools {
	return &QueryTools{querySvc: qSvc}
}

// ExecuteReadOnlyQuery validates the SQL statement with Guard before execution.
func (t *QueryTools) ExecuteReadOnlyQuery(ctx context.Context, conn connection.Connector, connID, connName, sqlQuery string) (*query.QueryResult, error) {
	if err := ValidateReadOnlyQuery(sqlQuery); err != nil {
		return nil, fmt.Errorf("read-only guard violation: %w", err)
	}
	return t.querySvc.Execute(ctx, conn, connID, connName, sqlQuery)
}
