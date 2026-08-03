package mysql

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
	_ "github.com/go-sql-driver/mysql"
)

// Connector implements core/connection.Connector for MySQL.
type Connector struct {
	db     *sql.DB
	config connection.ConnectionConfig
}

// NewConnector creates an uninitialized MySQL connector.
func NewConnector() connection.Connector {
	return &Connector{}
}

func (c *Connector) Connect(ctx context.Context, cfg connection.ConnectionConfig) error {
	c.config = cfg

	// DSN: username:password@tcp(host:port)/dbname?parseTime=true
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sql.Open("mysql", dsn)
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

func (c *Connector) currentDatabase(schema string) string {
	if schema != "" && schema != "public" {
		return schema
	}
	if c.config.Database != "" {
		return c.config.Database
	}
	return "mysql"
}

func (c *Connector) ListDatabases(ctx context.Context) ([]explorer.Database, error) {
	rows, err := c.db.QueryContext(ctx, "SHOW DATABASES;")
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
	dbName := c.currentDatabase(database)
	return []explorer.Schema{{Name: dbName}}, nil
}

func (c *Connector) ListTables(ctx context.Context, schema string) ([]explorer.Table, error) {
	dbName := c.currentDatabase(schema)
	querySQL := `
	SELECT table_name, table_type, COALESCE(table_rows, 0), COALESCE(table_comment, '')
	FROM information_schema.tables 
	WHERE table_schema = ?
	ORDER BY table_name;
	`
	rows, err := c.db.QueryContext(ctx, querySQL, dbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []explorer.Table
	for rows.Next() {
		var name, tblType, comment string
		var rowCount int64
		if err := rows.Scan(&name, &tblType, &rowCount, &comment); err != nil {
			return nil, err
		}
		tables = append(tables, explorer.Table{
			Name:     name,
			Schema:   dbName,
			Type:     tblType,
			RowCount: rowCount,
			Comment:  comment,
		})
	}
	return tables, rows.Err()
}

func (c *Connector) ListColumns(ctx context.Context, schema, table string) ([]explorer.Column, error) {
	dbName := c.currentDatabase(schema)
	querySQL := `
	SELECT 
		column_name, 
		column_type, 
		(is_nullable = 'YES'), 
		column_default, 
		ordinal_position, 
		(column_key = 'PRI'),
		COALESCE(extra, ''),
		COALESCE(column_comment, '')
	FROM information_schema.columns 
	WHERE table_schema = ? AND table_name = ?
	ORDER BY ordinal_position;
	`
	rows, err := c.db.QueryContext(ctx, querySQL, dbName, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []explorer.Column
	for rows.Next() {
		var name, cType, extra, comment string
		var isNullable, isPk bool
		var defVal sql.NullString
		var pos int

		if err := rows.Scan(&name, &cType, &isNullable, &defVal, &pos, &isPk, &extra, &comment); err != nil {
			return nil, err
		}

		var dPtr *string
		if defVal.Valid {
			v := defVal.String
			dPtr = &v
		}

		columns = append(columns, explorer.Column{
			Name:         name,
			DataType:     cType,
			IsNullable:   isNullable,
			IsPrimaryKey: isPk,
			DefaultValue: dPtr,
			Comment:      comment,
			Extra:        extra,
			Position:     pos,
		})
	}
	return columns, rows.Err()
}

func (c *Connector) ListIndexes(ctx context.Context, schema, table string) ([]explorer.Index, error) {
	dbName := c.currentDatabase(schema)
	querySQL := `
	SELECT 
		index_name, 
		non_unique = 0 as is_unique, 
		index_name = 'PRIMARY' as is_primary, 
		index_type, 
		column_name
	FROM information_schema.statistics
	WHERE table_schema = ? AND table_name = ?
	ORDER BY seq_in_index;
	`
	rows, err := c.db.QueryContext(ctx, querySQL, dbName, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	idxMap := make(map[string]*explorer.Index)
	for rows.Next() {
		var name, iType, colName string
		var isUnique, isPrimary bool
		if err := rows.Scan(&name, &isUnique, &isPrimary, &iType, &colName); err != nil {
			return nil, err
		}
		if existing, ok := idxMap[name]; ok {
			existing.Columns = append(existing.Columns, colName)
		} else {
			idxMap[name] = &explorer.Index{
				Name:      name,
				Columns:   []string{colName},
				IsUnique:  isUnique,
				IsPrimary: isPrimary,
				Type:      iType,
			}
		}
	}

	var indexes []explorer.Index
	for _, idx := range idxMap {
		indexes = append(indexes, *idx)
	}
	return indexes, rows.Err()
}

func (c *Connector) ListForeignKeys(ctx context.Context, schema, table string) ([]explorer.ForeignKey, error) {
	dbName := c.currentDatabase(schema)
	querySQL := `
	SELECT 
		k.constraint_name, 
		k.column_name, 
		k.referenced_table_name, 
		k.referenced_column_name,
		COALESCE(k.referenced_table_schema, ''),
		COALESCE(r.update_rule, 'RESTRICT'),
		COALESCE(r.delete_rule, 'RESTRICT')
	FROM information_schema.key_column_usage k
	LEFT JOIN information_schema.referential_constraints r
		ON r.constraint_name = k.constraint_name
		AND r.constraint_schema = k.table_schema
		AND r.table_name = k.table_name
	WHERE k.table_schema = ? 
	  AND k.table_name = ? 
	  AND k.referenced_table_name IS NOT NULL
	ORDER BY k.ordinal_position;
	`
	rows, err := c.db.QueryContext(ctx, querySQL, dbName, table)
	if err != nil {
		return nil, err
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
	dbName := c.currentDatabase(schema)
	querySQL := `
	SELECT table_name, view_definition 
	FROM information_schema.views 
	WHERE table_schema = ?;
	`
	rows, err := c.db.QueryContext(ctx, querySQL, dbName)
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
			Schema:     dbName,
			Definition: def,
		})
	}
	return views, rows.Err()
}

func (c *Connector) Execute(ctx context.Context, sqlQuery string) (*query.QueryResult, error) {
	return sqlutil.ExecuteGeneric(ctx, c.db, sqlQuery)
}

func (c *Connector) BrowseRows(ctx context.Context, schema, table string, opts data.BrowseOptions) (*data.RowPage, error) {
	dbName := c.currentDatabase(schema)
	placeholderFn := func(i int) string { return "?" }
	selectSQL, countSQL, args, countArgs := sqlutil.BuildBrowseQuery("`", dbName, table, opts, placeholderFn)

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

	colInfo, _ := c.ListColumns(ctx, dbName, table)
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
	dbName := c.currentDatabase(schema)
	var cols []string
	var placeholders []string
	var args []any

	for k, v := range values {
		cols = append(cols, fmt.Sprintf("`%s`", k))
		placeholders = append(placeholders, "?")
		args = append(args, v)
	}

	stmt := fmt.Sprintf("INSERT INTO `%s`.`%s` (%s) VALUES (%s)", dbName, table, strings.Join(cols, ", "), strings.Join(placeholders, ", "))
	_, err := c.db.ExecContext(ctx, stmt, args...)
	return err
}

func (c *Connector) UpdateRow(ctx context.Context, schema, table string, pk map[string]any, values map[string]any) error {
	dbName := c.currentDatabase(schema)
	var setClauses []string
	var args []any

	for k, v := range values {
		setClauses = append(setClauses, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}

	var whereClauses []string
	for k, v := range pk {
		whereClauses = append(whereClauses, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}

	stmt := fmt.Sprintf("UPDATE `%s`.`%s` SET %s WHERE %s", dbName, table, strings.Join(setClauses, ", "), strings.Join(whereClauses, " AND "))
	_, err := c.db.ExecContext(ctx, stmt, args...)
	return err
}

func (c *Connector) DeleteRow(ctx context.Context, schema, table string, pk map[string]any) error {
	dbName := c.currentDatabase(schema)
	var whereClauses []string
	var args []any

	for k, v := range pk {
		whereClauses = append(whereClauses, fmt.Sprintf("`%s` = ?", k))
		args = append(args, v)
	}

	stmt := fmt.Sprintf("DELETE FROM `%s`.`%s` WHERE %s", dbName, table, strings.Join(whereClauses, " AND "))
	_, err := c.db.ExecContext(ctx, stmt, args...)
	return err
}
