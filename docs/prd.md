# Product Requirements Document (PRD) — DBMW v0.1.0

- **Product Name:** DBMW (Database Management Workspace)
- **Version:** v0.1.0
- **Status:** Complete / Active

---

## 1. Executive Summary

DBMW is a fast, lightweight, self-hosted, local-first database management workspace. It ships as a single executable binary containing both a Go HTTP server and an embedded Vue 3 SPA frontend, launched from the CLI (`dbmw`).

It provides a developer-friendly alternative to bulky desktop database clients (like DBeaver, DataGrip, or pgAdmin) while keeping all database connection credentials and data strictly on the user's local machine.

---

## 2. Product Goals

1. **Single Binary Distribution**: Zero runtime dependencies (no Node.js required at runtime, no CGO compilation constraints).
2. **Sub-Second Boot Time**: Fast startup time (<1000ms) with low idle memory (<100MB RAM).
3. **Multi-Engine Consistency**: Unified user experience across PostgreSQL, MySQL, MariaDB, and SQLite.
4. **Interactive Schema Explorer**: Introspect databases, schemas, tables, columns, data types, primary keys, indexes, foreign keys, and views.
5. **Spreadsheet Data Grid**: Paginated data browsing, column sorting, custom filter rules, free-text search, inline editable cells, insert records, and delete rows.
6. **SQL Query Console**: SQL editor with `Ctrl+Enter` execution shortcut, execution timing badge, query history logging, and export to CSV/JSON.
7. **Entity-Relationship Diagram (ERD)**: Automatic visual diagram generation with table relationships mapped and zoom/pan/drag canvas.
8. **Framework Detection**: Automatic project scanner for Laravel, Prisma, Rails, Django, Go, and `dbmw.yml` configuration generator.
9. **AI Model Context Protocol (MCP)**: Stdio-based read-only tool interface for language models and AI agents.

---

## 3. Explicit Non-Goals (Out of Scope for v0.1.0)

- Multi-tenant user accounts and authentication.
- Cloud storage synchronization or telemetry tracking.
- Database server administration (e.g. creating physical database instances, managing system users/roles).
- AI query generation chat interfaces in the frontend (handled via external MCP clients).
