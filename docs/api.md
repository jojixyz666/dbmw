# REST API Reference — DBMW

All backend endpoints are prefixed with `/api`.

---

## 1. Health & Configuration

### `GET /api/health`
Returns system status.
- **Response `200 OK`**:
```json
{
  "status": "ok",
  "version": "0.1.0"
}
```

### `GET /api/config`
Fetches global workspace preferences.
- **Response `200 OK`**:
```json
{
  "theme": "dark",
  "defaultPageSize": 25,
  "serverPort": 8085,
  "autoOpenBrowser": true
}
```

### `POST /api/config`
Updates workspace preferences.
- **Request Body**: `AppConfig` JSON.

---

## 2. Connection Profiles

### `GET /api/connections`
Lists all saved connection profiles and the active connection ID.

### `POST /api/connections`
Creates or updates a connection configuration profile.
- **Request Body**:
```json
{
  "name": "Local PostgreSQL",
  "driver": "postgres",
  "host": "127.0.0.1",
  "port": 5432,
  "user": "postgres",
  "password": "secretpassword",
  "database": "dev_db",
  "sslMode": "disable"
}
```

### `POST /api/connections/test`
Tests whether a database connection configuration is reachable.

### `POST /api/connections/active`
Activates and switches current focus to a specified connection profile ID.
- **Request Body**: `{"id": "conn_123"}`

### `DELETE /api/connections/:id`
Deletes a saved database profile and terminates any active connection poolers.

---

## 3. Schema Explorer

- `GET /api/explorer/databases?connId=...` — Lists catalogues/databases.
- `GET /api/explorer/schemas?connId=...&database=...` — Lists database schemas.
- `GET /api/explorer/tables?connId=...&schema=...` — Lists tables in a schema.
- `GET /api/explorer/tables/:table/details?connId=...&schema=...` — Returns full columns, indexes, and FK constraints.
- `GET /api/explorer/columns/:table` — Returns columns information for a table.
- `GET /api/explorer/indexes/:table` — Returns index metadata for a table.
- `GET /api/explorer/foreign-keys/:table` — Returns foreign keys for a table.
- `GET /api/explorer/views` — Returns view definitions.

---

## 4. Query Execution

### `POST /api/query/execute`
Executes an arbitrary SQL query on the active or specified connection.
- **Request Body**:
```json
{
  "connectionId": "conn_123",
  "query": "SELECT * FROM users WHERE active = true LIMIT 10;"
}
```
- **Response `200 OK`**:
```json
{
  "columns": ["id", "username", "email", "created_at"],
  "columnTypes": ["INTEGER", "VARCHAR", "VARCHAR", "DATETIME"],
  "rows": [
    { "id": 1, "username": "alice", "email": "alice@example.com", "created_at": "2026-08-01T00:00:00Z" }
  ],
  "rowsAffected": 1,
  "executionTimeMs": 4.12
}
```

### `GET /api/query/history?connId=...&limit=50`
Returns recent historical queries with timestamps and durations.

### `POST /api/query/export/csv` & `POST /api/query/export/json`
Downloads query result datasets in CSV or JSON format.

---

## 5. Data Grid Manipulation

- `POST /api/data/browse/:table` — Browse paginated data with custom filters and column sorting.
- `POST /api/data/insert/:table` — Insert a new row into the table.
- `POST /api/data/update/:table` — Update a row matched by primary key.
- `POST /api/data/delete/:table` — Delete a row matched by primary key.

---

## 6. Visual ERD

### `GET /api/erd/generate?connId=...&schema=...`
Returns table nodes (with layout coordinates & column structures) and foreign key edges.

---

## 7. Project Detection

- `GET /api/project/detect?path=.` — Detects project framework and suggested database configuration.
- `POST /api/project/generate` — Writes `dbmw.yml` to the directory.
