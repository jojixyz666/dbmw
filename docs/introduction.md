# Introduction to DBMW

**DBMW (Database Management Workspace)** is a lightweight, local-first database management workspace designed as a modern, self-hosted alternative to heavy GUI database clients. It compiles into a **single static binary** that launches a local web server and serves an embedded Vue 3 SPA directly from your terminal.

---

## 🎯 Vision and Design Philosophy

Many traditional desktop database clients (e.g. DBeaver, pgAdmin, DataGrip) consume hundreds of megabytes of RAM, take seconds to start, and often require heavyweight runtimes like the JVM or Electron.

DBMW is built with three core principles:
1. **Single Binary & Zero External Runtime**: Ships as a single standalone executable per OS. The frontend is embedded via Go's `//go:embed`—no Node.js, Electron, or CGO dependencies required on the host system.
2. **Instant Boot (<1s) & Low Resource Footprint (<100MB RAM)**: Built on Go and lightweight Fiber/fasthttp to ensure immediate startup and minimal memory consumption during idle usage.
3. **Local-First & Private**: All connection credentials, query histories, and local configs live exclusively on your local machine (`~/.dbmw/`). No telemetry, no cloud syncing, and no external tracking.

---

## 🌟 Key Features

- **🔌 Supported Database Engines**:
  - **PostgreSQL** — High performance wire protocol via `pgx/v5`.
  - **MySQL & MariaDB** — Efficient native driver via `go-sql-driver/mysql`.
  - **SQLite** — Pure Go SQLite driver (`modernc.org/sqlite`) without CGO dependencies.
- **🗂️ Interactive Schema Explorer**: Browse databases, schemas, tables, views, column data types, nullability, defaults, primary keys, indexes, and foreign keys.
- **📋 Spreadsheet Data Grid**: Paginated data browsing with sortable columns, multi-rule filtering (`eq`, `neq`, `contains`, `gt`, `lt`, `is_null`, etc.), free-text search, inline editable cells, new record insertion, and row deletion.
- **⚡ SQL Query Console**: Fast multi-query SQL editor with `Ctrl+Enter` execution shortcut, elapsed execution timing (ms), historical query drawer with one-click reload, and dataset export to CSV and JSON.
- **🕸️ Visual ERD (Entity-Relationship Diagram)**: Automatic schema relationship graph mapping foreign keys into table cards with canvas zoom, pan, and dragging.
- **📦 Project Framework Detector**: Scans your active directory for frameworks (Laravel, Prisma/Node, Rails, Django, Go) and auto-generates `dbmw.yml`.
- **🤖 Read-Only Model Context Protocol (MCP)**: Native stdio-based MCP server guarded by an AST validator (`mcp/guard.go`) allowing AI assistants (Claude Desktop, Cursor, OpenCode) to introspect database structures safely.
