# Architecture Guide — DBMW v0.0.1

DBMW follows **Hexagonal Architecture (Ports and Adapters)** combined with a lightweight **Domain-Driven Design (DDD-Lite)** organization.

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

---

## 1. Architectural Layers & Separation of Concerns

### A. Pure Domain Layer (`core/`)
The core domain packages hold business entities, validation logic, and port interfaces. They have **zero dependencies** on external frameworks, SQL driver implementations, HTTP routers, or Vue components.

- `core/connection/`: Connection domain model, configuration validation, connection manager service, and the primary `Connector` & `ConnectionRepository` port interfaces.
- `core/explorer/`: Schema, table, column, index, view, and foreign key metadata structures.
- `core/query/`: SQL execution orchestration, timing metrics, query history tracking, and CSV/JSON formatting.
- `core/data/`: Data pagination, sorting, structured filter clauses, and CRUD manipulation logic.
- `core/erd/`: Layout calculation and relationship graph generation derived from table metadata and foreign keys.
- `core/project/`: Codebase inspection, environment variables parsing, framework detection, and `dbmw.yml` generation.

---

### B. Connector Adapters (`connector/`)
The connector layer provides concrete adapters for database engines, all satisfying the unified `connection.Connector` port:

- `connector/postgres/`: Native PostgreSQL adapter using `pgx/v5`.
- `connector/mysql/`: MySQL adapter using `go-sql-driver/mysql`.
- `connector/mariadb/`: MariaDB adapter with compatibility hooks over the MySQL protocol.
- `connector/sqlite/`: Modern pure Go SQLite adapter (`modernc.org/sqlite`).
- `connector/sqlutil/`: Shared dynamic row scanner (`*sql.Rows` -> `QueryResult`) and parameterized query builder.
- `connector/factory.go`: Centralized factory resolving `DriverType` to driver instances.

---

### C. Persistence Adapters (`storage/`)
Manages local state on the user's filesystem under `~/.dbmw/`:
- `storage/config_store.go`: Application settings (`config.json`).
- `storage/connection_store.go`: Saved database connection profiles (`connections.json`).
- `storage/history_store.go`: SQLite-backed historical query execution logs (`history.db`).

---

### D. Interface Adapters (`cli/`, `web/`, `mcp/`)
1. **CLI Layer (`cli/`)**: Cobra commands (`dbmw`, `open`, `connect`, `version`, `doctor`, `mcp`).
2. **Web Layer (`web/`)**: Fiber HTTP backend routing REST API requests to domain services and serving the embedded SPA via `//go:embed`.
3. **MCP Layer (`mcp/`)**: Protocol handler implementing Model Context Protocol over stdio for AI integration, strictly validated by `mcp/guard.go`.

---

## 2. Inviolable Architectural Principles

1. **Connector Port Immutability**: All database engines implement the same interface. Engine-specific extensions must remain internal to that engine adapter.
2. **Read-Only MCP Guard**: Any execution originating from `mcp/` must pass AST-level read-only verification before reaching the database connector.
3. **Zero Runtime Dependencies**: The compiled binary embeds all frontend static assets and uses pure-Go database drivers, requiring no external packages at runtime.
