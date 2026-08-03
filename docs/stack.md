# Technology Stack — DBMW v0.1.0

Every technical decision in DBMW is optimized for simplicity, performance, memory efficiency, and local-first reliability.

---

## 1. Backend

| Layer / Role | Selected Stack | Justification |
| :--- | :--- | :--- |
| **Language & Runtime** | **Go (1.24+)** | Single static binary per OS, low memory consumption (<100MB idle), native `//go:embed` asset support. |
| **Web Framework** | **Fiber v2 (`fasthttp`)** | Lightweight, high-throughput, low latency, Express-like developer ergonomics. |
| **CLI Framework** | **Cobra v1.10** | Standard Go CLI abstraction with structured subcommands, flags, and help menus. |
| **SQLite Persistence** | **`modernc.org/sqlite`** | 100% pure Go SQLite implementation without CGO compiler dependencies. |
| **PostgreSQL Driver** | **`pgx/v5` via `stdlib`** | Industrial-grade PostgreSQL wire protocol driver with connection pooling. |
| **MySQL / MariaDB Driver** | **`go-sql-driver/mysql`** | Efficient, battle-tested MySQL protocol driver. |

---

## 2. Frontend

| Layer / Role | Selected Stack | Justification |
| :--- | :--- | :--- |
| **UI Framework** | **Vue 3 (Composition API)** | Reactive, minimal overhead, easily bundled into static assets. |
| **Build Tooling** | **Vite 6** | Instant HMR development and minified Rollup production bundling (~49 KB gzipped). |
| **CSS Framework** | **Tailwind CSS v3** | Utility-first, midnight/emerald dark theme. |
| **State Management** | **Pinia v2** | Modular, type-safe application stores. |
| **Icons** | **Lucide Icons** | Clean SVG icon set. |

---

## 3. Storage Hierarchy

- Config file: `~/.dbmw/config.json`
- Connection profiles: `~/.dbmw/connections.json`
- Query execution history: `~/.dbmw/history.db`
- Optional project config: `<project_root>/dbmw.yml`
