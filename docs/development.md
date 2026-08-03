# Development & Extension Guide — DBMW

This guide explains how to develop, test, build, and extend DBMW.

---

## 1. Prerequisites

- **Go 1.24+**
- **Node.js 20+** & **npm 10+** (for frontend development and building assets)

---

## 2. Local Development Workflow

Start the backend and frontend separately during active development:

### Terminal 1: Backend
```bash
go run . --port 8085
```

### Terminal 2: Frontend (Vite Dev Server with HMR)
```bash
cd frontend
npm install
npm run dev
```
The Vite development server runs at `http://localhost:5173` and automatically proxies `/api` calls to `http://127.0.0.1:8085`.

---

## 3. Building for Production

Compile the Vue 3 SPA and embed it into the Go binary:

```bash
# 1. Build frontend dist
cd frontend
npm run build
cd ..

# 2. Compile static binary
go build -o dbmw.exe .
```

---

## 4. Running Test Suites

```bash
# Run all unit and integration tests across all packages
go test -v ./...

# Run tests with race condition detector
go test -race ./...
```

---

## 5. Adding a New Database Connector

To add support for a new database engine (e.g. `connector/clickhouse`, `connector/duckdb`):

1. Create a new directory under `connector/<engine>/`.
2. Implement the `connection.Connector` interface declared in `core/connection/port.go`:
   - `Connect(ctx, cfg)`, `Close()`, `Ping(ctx)`
   - `ListDatabases`, `ListSchemas`, `ListTables`, `ListColumns`, `ListIndexes`, `ListForeignKeys`, `ListViews`
   - `Execute(ctx, sqlQuery)`
   - `BrowseRows(ctx, schema, table, opts)`
   - `InsertRow`, `UpdateRow`, `DeleteRow`
3. Register your connector in `connector/factory.go`:
   ```go
   case connection.DriverClickHouse:
       return clickhouse.NewConnector(), nil
   ```
4. Add driver selection option to the frontend `ConnectionsModal.vue`.
5. Add unit tests for your driver under `connector/<engine>/`.
