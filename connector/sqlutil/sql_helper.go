package sqlutil

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"dbmw/core/data"
	"dbmw/core/query"
)

// ScanSQLRows transforms *sql.Rows into a QueryResult.
func ScanSQLRows(rows *sql.Rows) (*query.QueryResult, error) {
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	cols := make([]string, len(colTypes))
	typeNames := make([]string, len(colTypes))
	for i, ct := range colTypes {
		cols[i] = ct.Name()
		typeNames[i] = ct.DatabaseTypeName()
	}

	var resultRows []map[string]any

	for rows.Next() {
		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]any)
		for i, colName := range cols {
			val := columnPointers[i].(*any)
			switch v := (*val).(type) {
			case []byte:
				m[colName] = string(v)
			case time.Time:
				m[colName] = v.Format(time.RFC3339)
			default:
				m[colName] = v
			}
		}
		resultRows = append(resultRows, m)
	}

	return &query.QueryResult{
		Columns:      cols,
		ColumnTypes:  typeNames,
		Rows:         resultRows,
		RowsAffected: int64(len(resultRows)),
	}, rows.Err()
}

// ExecuteGeneric runs an arbitrary SQL statement and returns QueryResult.
func ExecuteGeneric(ctx context.Context, db *sql.DB, sqlQuery string) (*query.QueryResult, error) {
	trimmed := strings.TrimSpace(sqlQuery)
	upper := strings.ToUpper(trimmed)

	if strings.HasPrefix(upper, "SELECT") ||
		strings.HasPrefix(upper, "SHOW") ||
		strings.HasPrefix(upper, "PRAGMA") ||
		strings.HasPrefix(upper, "EXPLAIN") ||
		strings.HasPrefix(upper, "DESCRIBE") ||
		strings.HasPrefix(upper, "WITH") {
		rows, err := db.QueryContext(ctx, sqlQuery)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return ScanSQLRows(rows)
	}

	res, err := db.ExecContext(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	rowsAffected, _ := res.RowsAffected()
	return &query.QueryResult{
		Columns:      []string{"status", "rows_affected"},
		ColumnTypes:  []string{"VARCHAR", "INTEGER"},
		Rows:         []map[string]any{{"status": "OK", "rows_affected": rowsAffected}},
		RowsAffected: rowsAffected,
	}, nil
}

// BuildBrowseQuery constructs SQL and args for pagination, filter, and sort.
func BuildBrowseQuery(quoteChar string, schema, table string, opts data.BrowseOptions, placeholderFn func(int) string) (string, string, []any, []any) {
	var fullTable string
	if schema != "" && schema != "public" && schema != "main" {
		fullTable = fmt.Sprintf("%s%s%s.%s%s%s", quoteChar, schema, quoteChar, quoteChar, table, quoteChar)
	} else {
		fullTable = fmt.Sprintf("%s%s%s", quoteChar, table, quoteChar)
	}

	var whereClauses []string
	var args []any
	argIndex := 1

	for _, f := range opts.Filters {
		colQuoted := fmt.Sprintf("%s%s%s", quoteChar, f.Column, quoteChar)
		switch f.Operator {
		case data.OpEq:
			whereClauses = append(whereClauses, fmt.Sprintf("%s = %s", colQuoted, placeholderFn(argIndex)))
			args = append(args, f.Value)
			argIndex++
		case data.OpNeq:
			whereClauses = append(whereClauses, fmt.Sprintf("%s != %s", colQuoted, placeholderFn(argIndex)))
			args = append(args, f.Value)
			argIndex++
		case data.OpContains:
			whereClauses = append(whereClauses, fmt.Sprintf("%s LIKE %s", colQuoted, placeholderFn(argIndex)))
			args = append(args, fmt.Sprintf("%%%v%%", f.Value))
			argIndex++
		case data.OpStartsWith:
			whereClauses = append(whereClauses, fmt.Sprintf("%s LIKE %s", colQuoted, placeholderFn(argIndex)))
			args = append(args, fmt.Sprintf("%v%%", f.Value))
			argIndex++
		case data.OpEndsWith:
			whereClauses = append(whereClauses, fmt.Sprintf("%s LIKE %s", colQuoted, placeholderFn(argIndex)))
			args = append(args, fmt.Sprintf("%%%v", f.Value))
			argIndex++
		case data.OpGt:
			whereClauses = append(whereClauses, fmt.Sprintf("%s > %s", colQuoted, placeholderFn(argIndex)))
			args = append(args, f.Value)
			argIndex++
		case data.OpGte:
			whereClauses = append(whereClauses, fmt.Sprintf("%s >= %s", colQuoted, placeholderFn(argIndex)))
			args = append(args, f.Value)
			argIndex++
		case data.OpLt:
			whereClauses = append(whereClauses, fmt.Sprintf("%s < %s", colQuoted, placeholderFn(argIndex)))
			args = append(args, f.Value)
			argIndex++
		case data.OpLte:
			whereClauses = append(whereClauses, fmt.Sprintf("%s <= %s", colQuoted, placeholderFn(argIndex)))
			args = append(args, f.Value)
			argIndex++
		case data.OpIsNull:
			whereClauses = append(whereClauses, fmt.Sprintf("%s IS NULL", colQuoted))
		case data.OpIsNotNull:
			whereClauses = append(whereClauses, fmt.Sprintf("%s IS NOT NULL", colQuoted))
		}
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM %s%s", fullTable, whereSQL)
	countArgs := make([]any, len(args))
	copy(countArgs, args)

	orderSQL := ""
	if opts.SortBy != nil && opts.SortBy.Column != "" {
		dir := "ASC"
		if strings.ToLower(string(opts.SortBy.Direction)) == "desc" {
			dir = "DESC"
		}
		orderSQL = fmt.Sprintf(" ORDER BY %s%s%s %s", quoteChar, opts.SortBy.Column, quoteChar, dir)
	}

	pageSize := opts.PageSize
	if pageSize <= 0 {
		pageSize = 25
	}
	page := opts.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	limitSQL := fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
	selectSQL := fmt.Sprintf("SELECT * FROM %s%s%s%s", fullTable, whereSQL, orderSQL, limitSQL)

	return selectSQL, countSQL, args, countArgs
}
