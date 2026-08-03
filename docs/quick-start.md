# Quick Start Guide

This guide walks you through launching DBMW, connecting to a database, and using its core features.

---

## 1. Launching the Workspace

Start the DBMW server from anywhere in your terminal:

```bash
dbmw
```

This command will:
1. Initialize local configuration at `~/.dbmw/` (if running for the first time).
2. Start the local Fiber backend server (default: `http://127.0.0.1:8085`).
3. Automatically launch your default web browser to open the workspace UI.

> **Tip**: You can specify custom ports using `--port` or `-p` (e.g. `dbmw -p 9000`), or disable auto-opening the browser with `--no-open`.

---

## 2. Connecting to a Database

1. Click **+ Manage** or the connection selector in the top-left sidebar.
2. Choose your database engine:
   - **PostgreSQL**: Enter host (`127.0.0.1`), port (`5432`), database name, username, password, and SSL mode.
   - **MySQL / MariaDB**: Enter host (`127.0.0.1`), port (`3306`), database name, username, and password.
   - **SQLite**: Enter the path to your `.sqlite` or `.db` file (or `:memory:` for ephemeral storage).
3. Click **Test Connection** to verify database connectivity.
4. Click **Save Connection** to store the profile in `~/.dbmw/connections.json` and activate it.

---

## 3. Exploring Database Schema

Once connected, DBMW populates the sidebar with all discovered tables and views:
- **Click any Table** in the sidebar to open the **Schema Explorer** tab.
- Inspect detailed metadata:
  - Columns, data types, nullability, default values, and Primary Key indicators (🔑).
  - Defined Indexes and Unique constraints.
  - Foreign Key constraints mapping source columns to referenced tables.

---

## 4. Browsing and Editing Data (Spreadsheet Grid)

Switch to the **Data** tab on the navigation bar:
- **Pagination**: Navigate pages or adjust the page size (10, 25, 50, 100).
- **Sorting**: Click any column header to toggle ascending/descending order.
- **Search**: Type in the search box to filter records.
- **Filters**: Click **🎯 Filter** to add column criteria (e.g. `price > 100` or `status = 'active'`).
- **Inline Editing**: Click any cell, modify its value, and press `Enter` or click outside to instantly commit the update to the database.
- **Insert / Delete**: Use **+ Add Row** to insert new records or click 🗑️ on a row to delete it.

---

## 5. Executing Queries in the SQL Console

Switch to the **SQL** tab:
1. Write any SQL query in the editor (e.g. `SELECT * FROM users WHERE active = true;`).
2. Press **Ctrl+Enter** (or **Cmd+Enter** on macOS) or click **▶ Run Query**.
3. View the result table along with execution latency in milliseconds.
4. Click **📥 CSV** or **📥 JSON** to export the result dataset.
5. Click **📜 History** on the top action bar to inspect, review, or reload prior query runs.

---

## 6. Viewing Visual ERD

Switch to the **ERD** tab:
- DBMW automatically models all tables in the active schema as canvas cards.
- Foreign Key relationships are rendered as connecting SVG curves with source/target column endpoints.
- Use **＋ / －** to zoom, click and drag on the canvas background to **pan**, or drag individual table cards to reposition them.
