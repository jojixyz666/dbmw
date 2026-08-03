package mcp_test

import (
	"testing"
	"dbmw/mcp"
)

func TestMCPGuard(t *testing.T) {
	validQueries := []string{
		"SELECT * FROM users",
		"select id, name from products where price > 100",
		"EXPLAIN SELECT * FROM orders",
		"SHOW TABLES",
		"PRAGMA table_info('users')",
		"WITH active AS (SELECT id FROM users WHERE active = true) SELECT * FROM active",
		"DESCRIBE users",
		"-- Read only with comment prefix\nSELECT count(*) from items",
	}

	for _, q := range validQueries {
		if err := mcp.ValidateReadOnlyQuery(q); err != nil {
			t.Errorf("expected valid for %q, got error: %v", q, err)
		}
	}

	forbiddenQueries := []string{
		"INSERT INTO users (name) VALUES ('hacker')",
		"DROP TABLE secret",
		"UPDATE accounts SET balance = balance + 1000",
		"DELETE FROM logs",
		"ALTER TABLE users ADD COLUMN is_admin BOOLEAN",
		"TRUNCATE TABLE audit_log",
		"SELECT 1; DROP TABLE users;",
		"CREATE DATABASE exploit",
		"VACUUM",
	}

	for _, q := range forbiddenQueries {
		if err := mcp.ValidateReadOnlyQuery(q); err == nil {
			t.Errorf("expected error for forbidden write query %q, but got nil", q)
		}
	}
}
