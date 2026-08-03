package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"dbmw/connector/sqlutil"
	"dbmw/core/connection"
	"dbmw/core/data"
	"dbmw/core/explorer"
	"dbmw/core/query"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connector implements core/connection.Connector for PostgreSQL.
type Connector struct {
	db     *sql.DB
	config connection.ConnectionConfig
}

// NewConnector creates an uninitialized Postgres connector.
func NewConnector() connection.Connector {
	return &Connector{}
}

func (c *Connector) Connect(ctx context.Context, cfg connection.ConnectionConfig) error {
	c.config = cfg

	sslMode := cfg.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.Database,
		sslMode,
	)

	db, err := sql.Open("pgx", dsn)
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
	querySQL := `SELECT datname FROM pg_database WHERE datistemplate = false ORDER BY datname;`
	rows, err := c.db.QueryContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbs []explorer.Database
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		dbs = append(dbs, explorer.Database{Name: name})
	}
	return dbs, rows.Err()
}

func (c *Connector) ListSchemas(ctx context.Context, database string) ([]explorer.Schema, error) {
	querySQL := `
	SELECT schema_name 
	FROM information_schema.schemata 
	WHERE schema_name NOT IN ('information_schema', 'pg_catalog', 'pg_toast')
	ORDER BY schema_name;
	`
	rows, err := c.db.QueryContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []explorer.Schema
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		schemas = append(schemas, explorer.Schema{Name: name})
	}
	return schemas, rows.Err()
}

func (c *Connector) ListTables(ctx context.Context, schema string) ([]explorer.Table, error) {
	if schema == "" {
		schema = "public"
	}
	querySQL := `
	SELECT table_name, table_type 
	FROM information_schema.tables 
	WHERE table_schema = $1
	ORDER BY table_name;
	`
	rows, err := c.db.QueryContext(ctx, querySQL, schema)
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
			Schema: schema,
			Type:   tblType,
		})
	}
	return tables, rows.Err()
}

func (c *Connector) ListColumns(ctx context.Context, schema, table string) ([]explorer.Column, error) {
	if schema == "" {
		schema = "public"
	}
	querySQL := `
	SELECT 
		c.column_name, 
		c.data_type, 
		(c.is_nullable = 'YES') as is_nullable,
		c.column_default,
		c.ordinal_position,
		COALESCE(pk.is_pk, false) as is_pk
	FROM information_schema.columns c
	LEFT JOIN (
		SELECT kcu.column_name, true as is_pk
		FROM information_schema.table_constraints tc
		JOIN information_schema.key_column_usage kcu
			ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema
		WHERE tc.constraint_type = 'PRIMARY KEY' AND tc.table_schema = $1 AND tc.table_name = $2
	) pk ON c.column_name = pk.column_name
	WHERE c.table_schema = $1 AND c.table_name = $2
	ORDER BY c.ordinal_position;
	`
	rows, err := c.db.QueryContext(ctx, querySQL, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []explorer.Column
	for rows.Next() {
		var name, dType string
		var isNullable, isPk bool
		var defVal sql.NullString
		var pos int

		if err := rows.Scan(&name, &dType, &isNullable, &defVal, &pos, &isPk); err != nil {
			return nil, err
		}

		var dPtr *string
		if defVal.Valid {
			v := defVal.String
			dPtr = &v
		}

		columns = append(columns, explorer.Column{
			Name:         name,
			DataType:     dType,
			IsNullable:   isNullable,
			IsPrimaryKey: isPk,
			DefaultValue: dPtr,
			Position:     pos,
		})
	}
	return columns, rows.Err()
}

func (c *Connector) ListIndexes(ctx context.Context, schema, table string) ([]explorer.Index, error) {
	if schema == "" {
		schema = "public"
	}
	querySQL := `
	SELECT indexname, indexdef
	FROM pg_indexes
	WHERE schemaname = $1 AND tablename = $2;
	`
	rows, err := c.db.QueryContext(ctx, querySQL, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []explorer.Index
	for rows.Next() {
		var name, def string
		if err := rows.Scan(&name, &def); err != nil {
			return nil, err
		}
		isUnique := strings.Contains(strings.ToUpper(def), "UNIQUE INDEX")
		isPrimary := strings.HasSuffix(name, "_pkey") || strings.Contains(name, "pkey")
		indexes = append(indexes, explorer.Index{
			Name:      name,
			IsUnique:  isUnique,
			IsPrimary: isPrimary,
			Type:      "BTREE",
		})
	}
	return indexes, rows.Err()
}

func (c *Connector) ListForeignKeys(ctx context.Context, schema, table string) ([]explorer.ForeignKey, error) {
	if schema == "" {
		schema = "public"
	}
	querySQL := `
	SELECT
		tc.constraint_name,
		kcu.column_name,
		ccu.table_name AS foreign_table_name,
		ccu.column_name AS foreign_column_name,
		rc.update_rule,
		rc.delete_rule
	FROM information_schema.table_constraints AS tc
	JOIN information_schema.key_column_usage AS kcu
		ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema
	JOIN information_schema.constraint_column_usage AS ccu
		ON ccu.constraint_name = tc.constraint_name AND ccu.table_schema = tc.table_schema
	JOIN information_schema.referential_constraints AS rc
		ON rc.constraint_name = tc.constraint_name
	WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_schema = $1 AND tc.table_name = $2;
	`
	rows, err := c.db.QueryContext(ctx, querySQL, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fks []explorer.ForeignKey
	for rows.Next() {
		var fk explorer.ForeignKey
		if err := rows.Scan(&fk.Name, &fk.Column, &fk.RefTable, &fk.RefColumn, &fk.OnUpdate, &fk.OnDelete); err != nil {
			return nil, err
		}
		fks = append(fks, fk)
	}
	return fks, rows.Err()
}

func (c *Connector) ListViews(ctx context.Context, schema string) ([]explorer.View, error) {
	if schema == "" {
		schema = "public"
	}
	querySQL := `
	SELECT table_name, view_definition
	FROM information_schema.views
	WHERE table_schema = $1
	ORDER BY table_name;
	`
	rows, err := c.db.QueryContext(ctx, querySQL, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []explorer.View
	for rows.Next() {
		var name, def string
		if err := rows.Scan(&name, &def); err != nil {
			return nil, err
		}
		views = append(views, explorer.View{
			Name:       name,
			Schema:     schema,
			Definition: def,
		})
	}
	return views, rows.Err()
}

func (c *Connector) Execute(ctx context.Context, sqlQuery string) (*query.QueryResult, error) {
	return sqlutil.ExecuteGeneric(ctx, c.db, sqlQuery)
}

func (c *Connector) BrowseRows(ctx context.Context, schema, table string, opts data.BrowseOptions) (*data.RowPage, error) {
	if schema == "" {
		schema = "public"
	}
	placeholderFn := func(i int) string { return fmt.Sprintf("$%d", i) }
	selectSQL, countSQL, args, countArgs := sqlutil.BuildBrowseQuery(`"`, schema, table, opts, placeholderFn)

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
	if schema == "" {
		schema = "public"
	}
	var cols []string
	var placeholders []string
	var args []any
	idx := 1

	for k, v := range values {
		cols = append(cols, fmt.Sprintf(`"%s"`, k))
		placeholders = append(placeholders, fmt.Sprintf("$%d", idx))
		args = append(args, v)
		idx++
	}

	stmt := fmt.Sprintf(`INSERT INTO "%s"."%s" (%s) VALUES (%s)`, schema, table, strings.Join(cols, ", "), strings.Join(placeholders, ", "))
	_, err := c.db.ExecContext(ctx, stmt, args...)
	return err
}

func (c *Connector) UpdateRow(ctx context.Context, schema, table string, pk map[string]any, values map[string]any) error {
	if schema == "" {
		schema = "public"
	}
	var setClauses []string
	var args []any
	idx := 1

	for k, v := range values {
		setClauses = append(setClauses, fmt.Sprintf(`"%s" = $%d`, k, idx))
		args = append(args, v)
		idx++
	}

	var whereClauses []string
	for k, v := range pk {
		whereClauses = append(whereClauses, fmt.Sprintf(`"%s" = $%d`, k, idx))
		args = append(args, v)
		idx++
	}

	stmt := fmt.Sprintf(`UPDATE "%s"."%s" SET %s WHERE %s`, schema, table, strings.Join(setClauses, ", "), strings.Join(whereClauses, " AND "))
	_, err := c.db.ExecContext(ctx, stmt, args...)
	return err
}

func (c *Connector) DeleteRow(ctx context.Context, schema, table string, pk map[string]any) error {
	if schema == "" {
		schema = "public"
	}
	var whereClauses []string
	var args []any
	idx := 1

	for k, v := range pk {
		whereClauses = append(whereClauses, fmt.Sprintf(`"%s" = $%d`, k, idx))
		args = append(args, v)
		idx++
	}

	stmt := fmt.Sprintf(`DELETE FROM "%s"."%s" WHERE %s`, schema, table, strings.Join(whereClauses, " AND "))
	_, err := c.db.ExecContext(ctx, stmt, args...)
	return err
}
