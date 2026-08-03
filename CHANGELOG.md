# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.1] - 2026-08-03

### Added
- **Core Architecture & Hexagonal Ports**:
  - Implemented domain-per-folder pure architecture (`core/connection`, `core/explorer`, `core/query`, `core/data`, `core/erd`, `core/project`).
  - Standardized `Connector` interface contract across all database drivers.
  - Added comprehensive table-driven unit tests for all domain services.

- **Multi-Engine Database Connectors**:
  - **PostgreSQL**: Native wire driver adapter powered by `pgx/v5`.
  - **MySQL**: Connection and metadata adapter powered by `go-sql-driver/mysql`.
  - **MariaDB**: Compatible driver adapter wrapping MySQL protocol.
  - **SQLite**: Pure Go SQLite connector (`modernc.org/sqlite`) without CGO dependencies.

- **CLI Management Suite (Cobra)**:
  - `dbmw`: Starts local Fiber backend server and auto-launches workspace UI in default browser.
  - `dbmw open`: Reuses running server instance or starts a new one.
  - `dbmw connect`: Interactive terminal CLI for creating, listing, and testing database connection profiles.
  - `dbmw version`: Outputs detailed version metadata, Go runtime, and host OS/Arch.
  - `dbmw doctor`: Environment diagnostics verifying file permissions, `~/.dbmw` storage, port availability, and driver readiness.
  - `dbmw mcp`: Launches stdio-based Model Context Protocol server for AI assistant integration.

- **Read-Only Model Context Protocol (MCP)**:
  - `mcp/guard`: Strict read-only statement validator blocking all `INSERT`, `UPDATE`, `DELETE`, `DROP`, `ALTER`, `TRUNCATE`, `CREATE`, and write operations.
  - JSON-RPC 2.0 tools: `read_schemas`, `read_tables`, `read_columns`, and `execute_read_only_query`.

- **Storage Adapters**:
  - Local JSON storage for app configurations (`~/.dbmw/config.json`).
  - Local JSON storage for database profiles (`~/.dbmw/connections.json`).
  - Embedded SQLite history store for query logs and metrics (`~/.dbmw/history.db`).

- **Fiber Web API & Single Binary Embedding**:
  - High-performance RESTful API endpoints for backend communication.
  - Single-binary distribution with Vue 3 assets baked in via Go `//go:embed`.

- **Modern Vue 3 SPA UI**:
  - **Schema Explorer**: Interactive table structure browser for columns, types, primary keys, indexes, foreign keys, and views.
  - **Spreadsheet Data Grid**: Interactive table with pagination, column sorting, search, multi-condition filtering, inline cell editing, Add Row modal, and Delete Row confirmation.
  - **SQL Query Console**: SQL editor with `Ctrl+Enter` shortcut execution, execution time tracker (ms), query history drawer, export result to CSV and JSON.
  - **Visual ERD Canvas**: Entity-Relationship Diagram with automatic layout, Foreign Key mapping lines, zoom in/out, panning, and card dragging.
  - **Project Detector**: Codebase framework inspector (Laravel, Prisma, Rails, Django, Go) and `dbmw.yml` config generator.
  - **Connections & Settings Manager**: Modal for creation, editing, active switcher, and theme settings.
