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
	SELECT
		i.relname AS index_name,
		ix.indisunique AS is_unique,
		ix.indisprimary AS is_primary,
		am.amname AS index_type,
		COALESCE(string_agg(a.attname, ',' ORDER BY k.ord), '') AS columns
	FROM pg_catalog.pg_index ix
	JOIN pg_catalog.pg_class t ON t.oid = ix.indrelid
	JOIN pg_catalog.pg_class i ON i.oid = ix.indexrelid
	JOIN pg_catalog.pg_namespace ns ON ns.oid = t.relnamespace
	JOIN pg_catalog.pg_am am ON am.oid = i.relam
	JOIN unnest(ix.indkey) WITH ORDINALITY AS k(attnum, ord) ON true
	LEFT JOIN pg_catalog.pg_attribute a ON a.attrelid = t.oid AND a.attnum = k.attnum
	WHERE ns.nspname = $1 AND t.relname = $2
	GROUP BY i.relname, ix.indisunique, ix.indisprimary, am.amname
	ORDER BY i.relname;
	`
	rows, err := c.db.QueryContext(ctx, querySQL, schema, table)
	if err != nil {
		// Fallback to simpler pg_indexes query if complex catalog query fails
		fallbackSQL := `SELECT indexname, indexdef FROM pg_indexes WHERE schemaname = $1 AND tablename = $2;`
		fbRows, fbErr := c.db.QueryContext(ctx, fallbackSQL, schema, table)
		if fbErr != nil {
			return nil, err
		}
		defer fbRows.Close()

		var indexes []explorer.Index
		for fbRows.Next() {
			var name, def string
			if err := fbRows.Scan(&name, &def); err != nil {
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
		return indexes, fbRows.Err()
	}
	defer rows.Close()

	var indexes []explorer.Index
	for rows.Next() {
		var name, iType, colsStr string
		var isUnique, isPrimary bool
		if err := rows.Scan(&name, &isUnique, &isPrimary, &iType, &colsStr); err != nil {
			return nil, err
		}
		var cols []string
		if colsStr != "" {
			cols = strings.Split(colsStr, ",")
		}
		indexes = append(indexes, explorer.Index{
			Name:      name,
			Columns:   cols,
			IsUnique:  isUnique,
			IsPrimary: isPrimary,
			Type:      strings.ToUpper(iType),
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
		c.conname AS name,
		a.attname AS column_name,
		clf.relname AS foreign_table_name,
		af.attname AS foreign_column_name,
		COALESCE(nsf.nspname, '') AS foreign_schema_name,
		CASE c.confupdtype
			WHEN 'a' THEN 'NO ACTION'
			WHEN 'r' THEN 'RESTRICT'
			WHEN 'c' THEN 'CASCADE'
			WHEN 'n' THEN 'SET NULL'
			WHEN 'd' THEN 'SET DEFAULT'
			ELSE 'NO ACTION'
		END AS update_rule,
		CASE c.confdeltype
			WHEN 'a' THEN 'NO ACTION'
			WHEN 'r' THEN 'RESTRICT'
			WHEN 'c' THEN 'CASCADE'
			WHEN 'n' THEN 'SET NULL'
			WHEN 'd' THEN 'SET DEFAULT'
			ELSE 'NO ACTION'
		END AS delete_rule
	FROM pg_catalog.pg_constraint c
	JOIN pg_catalog.pg_class cl ON cl.oid = c.conrelid
	JOIN pg_catalog.pg_namespace ns ON ns.oid = cl.relnamespace
	JOIN pg_catalog.pg_class clf ON clf.oid = c.confrelid
	JOIN pg_catalog.pg_namespace nsf ON nsf.oid = clf.relnamespace
	JOIN unnest(c.conkey) WITH ORDINALITY AS k(attnum, ord) ON true
	JOIN pg_catalog.pg_attribute a ON a.attrelid = c.conrelid AND a.attnum = k.attnum
	JOIN unnest(c.confkey) WITH ORDINALITY AS f(attnum, ord) ON k.ord = f.ord
	JOIN pg_catalog.pg_attribute af ON af.attrelid = c.confrelid AND af.attnum = f.attnum
	WHERE c.contype = 'f'
	  AND ns.nspname = $1
	  AND cl.relname = $2
	ORDER BY c.conname, k.ord;
	`
	rows, err := c.db.QueryContext(ctx, querySQL, schema, table)
	if err != nil {
		// Fallback to information_schema if pg_catalog unnest isn't supported
		fallbackSQL := `
		SELECT
			tc.constraint_name,
			kcu.column_name,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name,
			COALESCE(rc.update_rule, 'NO ACTION'),
			COALESCE(rc.delete_rule, 'NO ACTION')
		FROM information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
			ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema
		JOIN information_schema.referential_constraints AS rc
			ON rc.constraint_name = tc.constraint_name AND rc.constraint_schema = tc.table_schema
		JOIN information_schema.constraint_column_usage AS ccu
			ON ccu.constraint_name = tc.constraint_name AND ccu.constraint_schema = tc.table_schema
		WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_schema = $1 AND tc.table_name = $2;
		`
		fbRows, fbErr := c.db.QueryContext(ctx, fallbackSQL, schema, table)
		if fbErr != nil {
			return nil, err
		}
		defer fbRows.Close()

		var fks []explorer.ForeignKey
		for fbRows.Next() {
			var fk explorer.ForeignKey
			if err := fbRows.Scan(&fk.Name, &fk.Column, &fk.RefTable, &fk.RefColumn, &fk.OnUpdate, &fk.OnDelete); err != nil {
				return nil, err
			}
			fks = append(fks, fk)
		}
		return fks, fbRows.Err()
	}
	defer rows.Close()

	var fks []explorer.ForeignKey
	for rows.Next() {
		var fk explorer.ForeignKey
		if err := rows.Scan(&fk.Name, &fk.Column, &fk.RefTable, &fk.RefColumn, &fk.RefSchema, &fk.OnUpdate, &fk.OnDelete); err != nil {
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
