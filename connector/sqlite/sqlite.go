package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"dbmw/connector/sqlutil"
	"dbmw/core/connection"
	"dbmw/core/data"
	"dbmw/core/explorer"
	"dbmw/core/query"
	_ "modernc.org/sqlite"
)

// Connector implements core/connection.Connector for SQLite.
type Connector struct {
	db *sql.DB
}

// NewConnector creates an uninitialized SQLite connector.
func NewConnector() connection.Connector {
	return &Connector{}
}

func (c *Connector) Connect(ctx context.Context, cfg connection.ConnectionConfig) error {
	path := cfg.FilePath
	if path == "" {
		path = cfg.Database
	}
	if path == "" {
		path = ":memory:"
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}
	c.db = db
	return c.Ping(ctx)
}

func (c *Connector) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *Connector) Ping(ctx context.Context) error {
	if c.db == nil {
		return fmt.Errorf("database is not connected")
	}
	return c.db.PingContext(ctx)
}

func (c *Connector) ListDatabases(ctx context.Context) ([]explorer.Database, error) {
	return []explorer.Database{{Name: "main"}}, nil
}

func (c *Connector) ListSchemas(ctx context.Context, database string) ([]explorer.Schema, error) {
	return []explorer.Schema{{Name: "main"}}, nil
}

func (c *Connector) ListTables(ctx context.Context, schema string) ([]explorer.Table, error) {
	querySQL := `
	SELECT name, type
	FROM sqlite_master 
	WHERE type IN ('table', 'view') AND name NOT LIKE 'sqlite_%'
	ORDER BY name;
	`
	rows, err := c.db.QueryContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []explorer.Table
	for rows.Next() {
		var name, tblType string
		if err := rows.Scan(&name, &tblType); err != nil {
			return nil, err
		}
		tables = append(tables, explorer.Table{
			Name:   name,
			Schema: "main",
			Type:   strings.ToUpper(tblType),
		})
	}
	return tables, rows.Err()
}

func (c *Connector) ListColumns(ctx context.Context, schema, table string) ([]explorer.Column, error) {
	querySQL := fmt.Sprintf(`PRAGMA table_info("%s");`, table)
	rows, err := c.db.QueryContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []explorer.Column
	var cid, notnull, pk int
	var name, colType string
	var dfltValue sql.NullString

	for rows.Next() {
		if err := rows.Scan(&cid, &name, &colType, &notnull, &dfltValue, &pk); err != nil {
			return nil, err
		}
		var defVal *string
		if dfltValue.Valid {
			v := dfltValue.String
			defVal = &v
		}

		columns = append(columns, explorer.Column{
			Name:         name,
			DataType:     colType,
			IsNullable:   notnull == 0,
			IsPrimaryKey: pk > 0,
			DefaultValue: defVal,
			Position:     cid + 1,
		})
	}

	return columns, rows.Err()
}

func (c *Connector) ListIndexes(ctx context.Context, schema, table string) ([]explorer.Index, error) {
	querySQL := fmt.Sprintf(`PRAGMA index_list("%s");`, table)
	rows, err := c.db.QueryContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []explorer.Index
	var seq, unique, partial int
	var name, origin string

	for rows.Next() {
		if err := rows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
			return nil, err
		}

		// Fetch columns in this index
		infoQuery := fmt.Sprintf(`PRAGMA index_info("%s");`, name)
		infoRows, err := c.db.QueryContext(ctx, infoQuery)
		var cols []string
		if err == nil {
			for infoRows.Next() {
				var seqno, cid int
				var colName string
				if err := infoRows.Scan(&seqno, &cid, &colName); err == nil {
					cols = append(cols, colName)
				}
			}
			infoRows.Close()
		}

		indexes = append(indexes, explorer.Index{
			Name:      name,
			Columns:   cols,
			IsUnique:  unique == 1,
			IsPrimary: origin == "pk",
			Type:      "BTREE",
		})
	}

	return indexes, rows.Err()
}

func (c *Connector) ListForeignKeys(ctx context.Context, schema, table string) ([]explorer.ForeignKey, error) {
	querySQL := fmt.Sprintf(`PRAGMA foreign_key_list("%s");`, table)
	rows, err := c.db.QueryContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fks []explorer.ForeignKey
	var id, seq int
	var refTable, fromCol, toCol, onUpdate, onDelete, match string

	for rows.Next() {
		if err := rows.Scan(&id, &seq, &refTable, &fromCol, &toCol, &onUpdate, &onDelete, &match); err != nil {
			return nil, err
		}
		fks = append(fks, explorer.ForeignKey{
			Name:      fmt.Sprintf("fk_%s_%s_%d", table, fromCol, id),
			Column:    fromCol,
			RefTable:  refTable,
			RefColumn: toCol,
			OnDelete:  onDelete,
			OnUpdate:  onUpdate,
		})
	}

	return fks, rows.Err()
}

func (c *Connector) ListViews(ctx context.Context, schema string) ([]explorer.View, error) {
	querySQL := `
	SELECT name, sql
	FROM sqlite_master
	WHERE type = 'view'
	ORDER BY name;
	`
	rows, err := c.db.QueryContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []explorer.View
	for rows.Next() {
		var name, definition string
		if err := rows.Scan(&name, &definition); err != nil {
			return nil, err
		}
		views = append(views, explorer.View{
			Name:       name,
			Schema:     "main",
			Definition: definition,
		})
	}
	return views, rows.Err()
}

func (c *Connector) Execute(ctx context.Context, sqlQuery string) (*query.QueryResult, error) {
	return sqlutil.ExecuteGeneric(ctx, c.db, sqlQuery)
}

func (c *Connector) BrowseRows(ctx context.Context, schema, table string, opts data.BrowseOptions) (*data.RowPage, error) {
	placeholderFn := func(i int) string { return "?" }
	selectSQL, countSQL, args, countArgs := sqlutil.BuildBrowseQuery(`"`, "", table, opts, placeholderFn)

	var totalRows int64
	if err := c.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&totalRows); err != nil {
		totalRows = 0
	}

	rows, err := c.db.QueryContext(ctx, selectSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	qRes, err := sqlutil.ScanSQLRows(rows)
	if err != nil {
		return nil, err
	}

	pageSize := opts.PageSize
	if pageSize <= 0 {
		pageSize = 25
	}
	totalPages := int((totalRows + int64(pageSize) - 1) / int64(pageSize))
	if totalPages < 1 {
		totalPages = 1
	}

	// Fetch primary key columns
	colInfo, _ := c.ListColumns(ctx, schema, table)
	var pkNames []string
	for _, col := range colInfo {
		if col.IsPrimaryKey {
			pkNames = append(pkNames, col.Name)
		}
	}

	return &data.RowPage{
		Page:       opts.Page,
		PageSize:   pageSize,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Columns:    qRes.Columns,
		Rows:       qRes.Rows,
		PrimaryKey: pkNames,
	}, nil
}

func (c *Connector) InsertRow(ctx context.Context, schema, table string, values map[string]any) error {
	var cols []string
	var placeholders []string
	var args []any

	for k, v := range values {
		cols = append(cols, fmt.Sprintf(`"%s"`, k))
		placeholders = append(placeholders, "?")
		args = append(args, v)
	}

	stmt := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`, table, strings.Join(cols, ", "), strings.Join(placeholders, ", "))
	_, err := c.db.ExecContext(ctx, stmt, args...)
	return err
}

func (c *Connector) UpdateRow(ctx context.Context, schema, table string, pk map[string]any, values map[string]any) error {
	var setClauses []string
	var args []any

	for k, v := range values {
		setClauses = append(setClauses, fmt.Sprintf(`"%s" = ?`, k))
		args = append(args, v)
	}

	var whereClauses []string
	for k, v := range pk {
		whereClauses = append(whereClauses, fmt.Sprintf(`"%s" = ?`, k))
		args = append(args, v)
	}

	stmt := fmt.Sprintf(`UPDATE "%s" SET %s WHERE %s`, table, strings.Join(setClauses, ", "), strings.Join(whereClauses, " AND "))
	_, err := c.db.ExecContext(ctx, stmt, args...)
	return err
}

func (c *Connector) DeleteRow(ctx context.Context, schema, table string, pk map[string]any) error {
	var whereClauses []string
	var args []any

	for k, v := range pk {
		whereClauses = append(whereClauses, fmt.Sprintf(`"%s" = ?`, k))
		args = append(args, v)
	}

	stmt := fmt.Sprintf(`DELETE FROM "%s" WHERE %s`, table, strings.Join(whereClauses, " AND "))
	_, err := c.db.ExecContext(ctx, stmt, args...)
	return err
}
