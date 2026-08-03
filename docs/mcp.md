# Model Context Protocol (MCP) Guide — DBMW

DBMW provides built-in support for the **Model Context Protocol (MCP)**, allowing AI language models and coding assistants (such as Claude Desktop, OpenCode, and Cursor) to introspect database structures and run analytical queries over stdio.

---

## 1. Safety Guard (`mcp/guard.go`)

To avoid unintended data loss, DBMW strictly enforces a **Read-Only Safety Guard**:
- **Allowed Query Statements**: `SELECT`, `EXPLAIN`, `SHOW`, `DESCRIBE`, `PRAGMA` (read-only), and `WITH ... SELECT`.
- **Forbidden Statement Keywords**: `INSERT`, `UPDATE`, `DELETE`, `DROP`, `ALTER`, `TRUNCATE`, `CREATE`, `REPLACE`, `GRANT`, `REVOKE`, `ATTACH`, `DETACH`, `BEGIN`, `COMMIT`, `ROLLBACK`, `VACUUM`.
- **Multi-Statement Guard**: Queries with multiple separated statements (`SELECT 1; DROP TABLE users;`) are strictly blocked before execution.

---

## 2. Available MCP Tools

| Tool Indicator | Purpose | Parameters |
| :--- | :--- | :--- |
| `read_schemas` | List all available schemas in database catalog. | `connectionId`, `database` |
| `read_tables` | List all tables and views in schema. | `connectionId`, `schema` |
| `read_columns` | Inspect column names, data types, and primary keys. | `connectionId`, `schema`, `table` |
| `execute_read_only_query` | Run a validated SELECT query against database. | `connectionId`, `query` |

---

## 3. Configuring Claude Desktop or AI Assistants

Add the following to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "dbmw": {
      "command": "D:/opensource/dbmw/dbmw.exe",
      "args": ["mcp"]
    }
  }
}
```
