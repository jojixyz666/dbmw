# Project Memory — DBMW (Database Management Workspace)

**Last Updated**: 2026-08-03
**Status**: Active / Production v0.0.1 Foundation Release

---

## 1. Project Overview & URLs
- **Repository**: https://github.com/jojixyz666/dbmw
- **Online Documentation Website**: https://jojixyz666.github.io/dbmw/
- **Description**: Lightweight, self-hosted, local-first database management workspace that compiles into a single static Go binary with an embedded Vue 3 SPA.

---

## 2. Tech Stack & Dependencies
- **Core / Backend**: Go 1.24+, Fiber v2 (`fasthttp`), Cobra CLI v1.10.
- **Database Engine Drivers**:
  - PostgreSQL: `github.com/jackc/pgx/v5`
  - MySQL / MariaDB: `github.com/go-sql-driver/mysql`
  - SQLite: `modernc.org/sqlite` (100% Pure Go, zero CGO requirement).
- **Persistent Storage**: `~/.dbmw/` (`config.json`, `connections.json`, `history.db` SQLite).
- **AI Integration**: Built-in Model Context Protocol (MCP) server over stdio with AST-level read-only safety guard (`mcp/guard.go`).
- **Frontend SPA**: Vue 3 (Composition API), Vite 6, Tailwind CSS v3, Pinia v2, embedded at compile-time via Go `//go:embed all:dist`.
- **Documentation Engine**: MkDocs with Material for MkDocs theme.

---

## 3. High-Priority Conventions & Rules
1. **Hexagonal Domain Purity**: Never import Fiber, database drivers, or frontend packages inside `core/<domain>/`. Only interfaces in `port.go` and pure Go standard library.
2. **One Port per Capability**: All database connectors must implement `core/connection.Connector`.
3. **MCP Read-Only Enforcement**: Any MCP query execution must pass validation in `mcp/guard.go`.
4. **Single Binary Embed**: The frontend must remain embedded inside binary; runtime must not depend on Node.js.

---

## 4. Key Workflows & Commands
- **Build Full Binary**:
  ```bash
  cd frontend && npm run build && cd ..
  go build -o dbmw.exe .
  ```
- **Run Tests**:
  ```bash
  go test -v ./...
  ```
- **CLI Commands**:
  - Run app: `./dbmw.exe` (starts server on port 8085 and opens browser)
  - Health check: `./dbmw.exe doctor`
  - Interactive connections: `./dbmw.exe connect`
  - Version info: `./dbmw.exe version`
  - Stdio MCP server: `./dbmw.exe mcp`
- **Deploy Documentation**:
  ```bash
  python -m mkdocs gh-deploy
  ```
