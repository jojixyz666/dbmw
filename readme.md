# DBMW (Database Management Workspace)

[![License: MIT](https://img.shields.io/badge/License-MIT-emerald.svg)](./LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-blue.svg)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.x-42b883.svg)](https://vuejs.org)

**DBMW** is a lightweight, local-first database management workspace that compiles into a single static binary. It provides a fast, self-hosted web UI launched directly from the terminal, keeping all your database credentials and data securely on your local machine.

---

## ✨ Key Features

- **🚀 Single Static Binary**: Embedded Vue 3 SPA via Go's `//go:embed` — zero runtime dependencies, no Node.js or CGO required.
- **🔌 Multi-Engine Database Support**:
  - PostgreSQL (via `pgx/v5`)
  - MySQL & MariaDB (via `go-sql-driver/mysql`)
  - SQLite (pure Go via `modernc.org/sqlite`)
- **🗂️ Interactive Schema Explorer**: Introspect databases, schemas, tables, columns, data types, primary keys, indexes, foreign keys, and views.
- **📋 Spreadsheet Data Grid**: Paginated data browsing with sorting, multi-condition filters, free-text search, inline cell editing, Add Row dialog, and Delete Row confirmation.
- **⚡ SQL Query Console**: Multi-query editor with `Ctrl+Enter` execution shortcut, execution timing badge, historical queries drawer, and export to CSV/JSON.
- **🕸️ Visual ERD (Entity-Relationship Diagram)**: Automatic diagram generator mapping foreign key relationships with zoom/pan and draggable table cards.
- **📦 Project Framework Detector**: Auto-detects project framework configurations (Laravel, Prisma, Rails, Django, Go) and generates `dbmw.yml`.
- **🤖 Built-in MCP Server**: Read-only Model Context Protocol server over stdio protected by a strict `guard.go` validator for AI assistants (Claude Desktop, Cursor, OpenCode).

---

## 🚀 Quick Start

### Build & Run

```bash
# Build the binary
go build -o dbmw.exe .

# Launch server + auto-open web UI
./dbmw.exe
```

### CLI Subcommands

```bash
# System environment & storage health checks
./dbmw.exe doctor

# Interactively test & manage connection profiles
./dbmw.exe connect

# Launch MCP stdio server for AI tools
./dbmw.exe mcp

# Print version and build telemetry
./dbmw.exe version
```

---

## 📖 Documentation

Detailed guides and documentation are available in the [`docs/`](./docs) directory:

- [**Introduction (`docs/introduction.md`)**](./docs/introduction.md) — Product overview, principles, and feature breakdown.
- [**Installation Guide (`docs/installation.md`)**](./docs/installation.md) — Pre-built binaries and building from source.
- [**Quick Start (`docs/quick-start.md`)**](./docs/quick-start.md) — Step-by-step walkthrough to get up and running.
- [**Architecture Guide (`docs/architecture.md`)**](./docs/architecture.md) — Hexagonal ports, domain isolation, and design decisions.
- [**Configuration Guide (`docs/configuration.md`)**](./docs/configuration.md) — Local settings, database profiles, and AI MCP integration.
- [**Development Guide (`docs/development.md`)**](./docs/development.md) — Contributing, testing, and adding new database connectors.

---

## 📜 License

[MIT License](./LICENSE) © 2026 M-Code
