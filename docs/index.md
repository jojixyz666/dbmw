# DBMW Documentation

<div align="center">

![DBMW Banner](https://img.shields.io/badge/DBMW-v0.0.1%20Foundation-10b981?style=for-the-badge)
[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue.js-3.x-42b883?style=for-the-badge&logo=vuedotjs&logoColor=white)](https://vuejs.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

<h3>Database Management Workspace</h3>

<p><strong>A lightweight, local-first database management workspace that compiles into a single static binary.</strong></p>

</div>

---

## ⚡ Overview

**DBMW** provides a fast, self-hosted, web-based database management UI launched directly from the terminal. All data and connection credentials remain 100% on your local machine, giving you total control and data privacy.

```bash
# Launch server and auto-open UI in your default browser
dbmw
```

---

## 🚀 Key Highlights

=== "🔌 Multi-Engine Support"
    Connect to **PostgreSQL**, **MySQL**, **MariaDB**, and **SQLite** through one unified, consistent workspace interface. SQLite uses a pure Go driver without CGO required.

=== "🗂️ Interactive Schema Explorer"
    Explore physical databases, schemas, tables, views, column data types, nullability, primary keys, indexes, and full foreign key constraints.

=== "📋 Spreadsheet Data Grid"
    Browse records with pagination, sort by column, filter by multi-criteria rules, perform free-text search, edit cell values inline, insert new rows, and delete records with confirmation.

=== "⚡ SQL Query Console"
    Multi-query editor with syntax highlighting, `Ctrl+Enter` shortcut execution, elapsed timing in milliseconds, historical query drawer with one-click re-run, and data export to CSV and JSON.

=== "🕸️ Visual ERD Diagram"
    Automatic Entity-Relationship Diagram generator directly mapping Foreign Key relationships into interactive table cards with zoom, pan, and dragging.

=== "🤖 Built-in MCP Server"
    A native Model Context Protocol (MCP) server running over stdio, guarded by an AST-level read-only validator for AI assistants (such as Claude Desktop, Cursor, and OpenCode).

---

## 📚 Documentation Sections

| Section | Description |
| :--- | :--- |
| [**Introduction**](introduction.md) | Vision, philosophy, and supported database drivers. |
| [**Installation**](installation.md) | How to download pre-built binaries or build from source. |
| [**Quick Start**](quick-start.md) | Connect to your first database, explore schemas, and run queries. |
| [**Architecture**](architecture.md) | Hexagonal Architecture (Ports and Adapters) & DDD-Lite design. |
| [**Configuration**](configuration.md) | Local settings (`~/.dbmw/`), project configs (`dbmw.yml`), and MCP setup. |
| [**Development**](development.md) | Contributing guide, local dev server, testing, and adding new database connectors. |

---

## 📦 Single Binary Distribution

The frontend Vue 3 single page application is compiled and embedded into the Go binary at compile time via Go's `//go:embed`. There is **zero runtime dependency** on Node.js or any external web server on target machines.
