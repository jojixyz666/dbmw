# CLI Reference — DBMW

DBMW provides a command-line interface based on **spf13/cobra**.

---

## Commands Catalog

### 1. `dbmw` (Root Command)
Starts the Fiber HTTP server and automatically opens the workspace UI in your default browser.

```bash
# Start on default port (8085)
./dbmw.exe

# Start on custom port
./dbmw.exe --port 9000
./dbmw.exe -p 9000

# Start server without auto-opening web browser
./dbmw.exe --no-open
```

---

### 2. `dbmw open`
Connects to an already-running DBMW server instance or starts a new instance if none exists.

```bash
./dbmw.exe open
```

---

### 3. `dbmw connect`
Interactively list, inspect, create, and test database connection profiles directly from your terminal.

```bash
./dbmw.exe connect
```

**Walkthrough**:
1. Choose to test an existing saved connection or press `[A]` to add a new profile.
2. Select driver: `postgres`, `mysql`, `mariadb`, or `sqlite`.
3. Provide connection details (host, port, credentials, or file path).
4. Automated reachability ping test runs immediately before saving profile to `~/.dbmw/connections.json`.

---

### 4. `dbmw doctor`
Runs automated diagnostic health checks on your local environment:
- Resolves and verifies `~/.dbmw/` storage directory permissions.
- Validates integrity of `config.json` and `connections.json`.
- Inspects SQLite query history database `history.db`.
- Tests TCP port binding availability.
- Verifies registration of all 4 database driver factories.

```bash
./dbmw.exe doctor
```

---

### 5. `dbmw mcp`
Launches the Model Context Protocol (MCP) server over stdin/stdout, allowing AI coding assistants (such as Claude Desktop or OpenCode agents) to query connected database schemas and run read-only queries.

```bash
./dbmw.exe mcp
```

---

### 6. `dbmw version`
Outputs build version tag, commit hash, build date, Go runtime, and system architecture.

```bash
./dbmw.exe version
```
