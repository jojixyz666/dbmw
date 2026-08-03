# Configuration Guide

DBMW stores all local preferences, database profiles, and query logs under your user directory in `~/.dbmw/`.

---

## 1. Global Application Settings (`~/.dbmw/config.json`)

Created automatically on first launch:

```json
{
  "theme": "dark",
  "defaultPageSize": 25,
  "serverPort": 8085,
  "autoOpenBrowser": true
}
```

### Configuration Options:
- `theme`: UI visual theme (`"dark"` or `"light"`).
- `defaultPageSize`: Default rows per page in the Spreadsheet Data Grid (`10`, `25`, `50`, `100`).
- `serverPort`: Default listening port for the web server (can be overridden via `--port` flag).
- `autoOpenBrowser`: Automatically launch default browser when running `dbmw` (`true` or `false`).

---

## 2. Saved Connection Profiles (`~/.dbmw/connections.json`)

Saved database profiles are persisted as structured JSON:

```json
[
  {
    "id": "conn_1738598400000000000",
    "name": "Local PostgreSQL Dev",
    "driver": "postgres",
    "host": "127.0.0.1",
    "port": 5432,
    "user": "postgres",
    "password": "secretpassword",
    "database": "app_development",
    "sslMode": "disable",
    "createdAt": "2026-08-03T10:00:00Z",
    "updatedAt": "2026-08-03T10:00:00Z"
  },
  {
    "id": "conn_1738598401000000001",
    "name": "Project SQLite",
    "driver": "sqlite",
    "filePath": "./database.sqlite",
    "createdAt": "2026-08-03T10:00:00Z",
    "updatedAt": "2026-08-03T10:00:00Z"
  }
]
```

---

## 3. Project Configuration (`dbmw.yml`)

DBMW supports project-level configuration files placed in the root of your project directory. This enables team members to share connection templates without committing sensitive passwords.

Example `dbmw.yml`:

```yaml
version: "1"
project_name: "My Application"
default_connection: "primary_dev"
connections:
  - name: "primary_dev"
    driver: "postgres"
    host: "127.0.0.1"
    port: 5432
    user: "postgres"
    database: "my_app_development"
  - name: "sqlite_test"
    driver: "sqlite"
    file_path: "./test.db"
```

---

## 4. MCP Server Setup for AI Assistants

You can integrate DBMW with AI assistants like Claude Desktop, Cursor, or OpenCode by registering the stdio MCP bridge.

### Claude Desktop (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "dbmw": {
      "command": "/usr/local/bin/dbmw",
      "args": ["mcp"]
    }
  }
}
```

On Windows:
```json
{
  "mcpServers": {
    "dbmw": {
      "command": "C:\\tools\\dbmw.exe",
      "args": ["mcp"]
    }
  }
}
```
