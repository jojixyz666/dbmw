package mcp

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrWriteNotAllowed     = errors.New("MCP is strictly read-only: write/mutation queries are prohibited")
	ErrMultipleStatements  = errors.New("MCP does not permit multiple queries in a single execution")
	ErrDisallowedStatement = errors.New("statement is not allowed in MCP read-only mode")
)

var forbiddenKeywords = []string{
	"INSERT", "UPDATE", "DELETE", "DROP", "ALTER", "TRUNCATE",
	"CREATE", "REPLACE", "GRANT", "REVOKE", "ATTACH", "DETACH",
	"BEGIN", "COMMIT", "ROLLBACK", "VACUUM", "REINDEX",
}

// ValidateReadOnlyQuery inspects an SQL query to strictly guarantee read-only execution.
func ValidateReadOnlyQuery(sqlQuery string) error {
	trimmed := strings.TrimSpace(sqlQuery)
	if trimmed == "" {
		return errors.New("query cannot be empty")
	}

	// Disallow multiple semicolon-separated statements
	statements := splitStatements(trimmed)
	if len(statements) > 1 {
		return ErrMultipleStatements
	}

	// Normalize tokens
	cleanSQL := removeComments(trimmed)
	words := extractWords(cleanSQL)
	if len(words) == 0 {
		return errors.New("empty query after stripping comments")
	}

	for _, w := range words {
		upperWord := strings.ToUpper(w)
		for _, forbidden := range forbiddenKeywords {
			if upperWord == forbidden {
				return ErrWriteNotAllowed
			}
		}
	}

	firstWord := strings.ToUpper(words[0])
	allowedPrefixes := map[string]bool{
		"SELECT":   true,
		"EXPLAIN":  true,
		"SHOW":     true,
		"PRAGMA":   true,
		"DESCRIBE": true,
		"DESC":     true,
		"WITH":     true,
	}

	if !allowedPrefixes[firstWord] {
		return ErrDisallowedStatement
	}

	return nil
}

func removeComments(s string) string {
	var result strings.Builder
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "--") {
			continue
		}
		if idx := strings.Index(l, "--"); idx != -1 {
			l = l[:idx]
		}
		result.WriteString(l)
		result.WriteString(" ")
	}
	return result.String()
}

func splitStatements(s string) []string {
	var results []string
	parts := strings.Split(s, ";")
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			results = append(results, strings.TrimSpace(p))
		}
	}
	return results
}

func extractWords(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
	})
}
