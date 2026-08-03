# Architecture — DBMW v0.1.0

This document outlines the architectural blueprint of DBMW.

---

## 1. Architectural Pattern: Hexagonal Architecture (Ports & Adapters)

DBMW is built with **Hexagonal Architecture with DDD-Lite domain-per-folder packaging**:

```text
┌────────────────────────────────────────────────────────┐
│                        ADAPTERS                        │
│                                                        │
│   CLI (Cobra)       Fiber HTTP (web/)       MCP (mcp/) │
│        │                   │                    │      │
│        ▼                   ▼                    ▼      │
│ ┌────────────────────────────────────────────────────┐ │
│ │                  CORE DOMAIN HOOKS                 │ │
│ │                                                    │ │
│ │  connection  explorer  query  data  erd  project   │ │
│ │                                                    │ │
│ └────────────────────────────────────────────────────┘ │
│        ▲                   ▲                    ▲      │
│        │                   │                    │      │
│   Postgres/MySQL     SQLite (Pure Go)   Storage Store  │
│                                                        │
│                    PORT IMPLEMENTATIONS                │
└────────────────────────────────────────────────────────┘
```

### Core Invariants:
1. **Domain Purity (`core/`)**: Domain services directly depend only on standard Go libraries and domain interfaces (`port.go`). They **never** import Fiber, specific database driver packages, or Vue/frontend constructs.
2. **One Port per Capability (`core/connection/port.go`)**: All database engine adapters implement the unified `Connector` interface.
3. **Strict Safety Guard (`mcp/guard.go`)**: All queries executed via the MCP protocol are strictly constrained to read-only queries with multi-level AST inspection.
4. **Embedded Single Artifact**: Vue 3 SPA is compiled into `frontend/dist` and embedded seamlessly via Go's `//go:embed`.

---

## 2. Directory Layout Breakdown

```
dbmw/
├── cli/                      # Cobra-based CLI entrypoints (driving adapters)
│   ├── root.go               # Server runner + launch options
│   ├── open.go               # Focus / startup handler
│   ├── connect.go            # Interactive terminal DB profile manager
│   ├── version.go            # Build & runtime version output
│   ├── doctor.go             # System health diagnostics
│   └── mcp.go                # Stdio MCP server runner
│
├── core/                     # Pure domain logic + ports (hexagon center)
│   ├── connection/           # Connection configs, lifecycle, & Connector port
│   ├── explorer/             # Database, schema, table, & index metadata
│   ├── query/                # SQL execution, timing metrics, and export
│   ├── data/                 # Spreadsheet row browsing, filters, & CRUD
│   ├── erd/                  # Automatic ERD relationship graph builder
│   └── project/              # Framework detection & dbmw.yml generator
│
├── connector/                # Engine driver adapters (driven ports)
│   ├── sqlutil/              # Shared sql.Rows parsing and SQL builders
│   ├── postgres/             # pgx/v5 PostgreSQL implementation
│   ├── mysql/                # go-sql-driver/mysql MySQL implementation
│   ├── mariadb/              # MariaDB protocol implementation
│   ├── sqlite/               # modernc.org/sqlite pure Go implementation
│   └── factory.go            # Driver factory resolver
│
├── storage/                  # Local persistence adapters (~/.dbmw/)
│   ├── paths.go              # ~/.dbmw resolution
│   ├── config_store.go       # ~/.dbmw/config.json
│   ├── connection_store.go   # ~/.dbmw/connections.json
│   └── history_store.go      # ~/.dbmw/history.db (SQLite logs)
│
├── web/                      # Driving adapter: Fiber web backend
│   ├── server.go             # Server initialization & lifecycle
│   ├── router.go             # Route registrations and middlewares
│   ├── embed.go              # Embedded frontend filesystem wrapper
│   └── handlers/             # REST API controllers
│
├── mcp/                      # Driving adapter: Read-only MCP server
│   ├── guard.go              # AST & keyword read-only enforcement
│   ├── server.go             # JSON-RPC 2.0 stdio loop
│   ├── tools_schema.go       # Schema exploration tools
│   └── tools_query.go        # Read-only query runner
│
└── frontend/                 # Vue 3 SPA application source
    ├── src/
    │   ├── App.vue           # Main application layout
    │   ├── stores/app.js     # Pinia state management
    │   └── components/       # Explorer, SqlEditor, DataGrid, ERD, Project
    └── dist/                 # Compiled assets for //go:embed
```

---

## 3. Data and Execution Flow

```text
[User Browser] ──HTTP──> [Fiber Router (web/)] ──> [Handlers]
                                                        │
                      ┌─────────────────────────────────┘
                      ▼
            [core/connection/service.go]
                      │
           ┌──────────┴──────────┐
           ▼                     ▼
    [Connector Port]      [Storage Repositories]
           │                     │
   ┌───────┼───────┐      (~/.dbmw/*.json, history.db)
   ▼       ▼       ▼
Postgres MySQL  SQLite
```
