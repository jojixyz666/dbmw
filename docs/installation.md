# Installation Guide

DBMW is distributed as a single static binary for all major operating systems.

---

## 1. System Requirements

- **Supported Platforms**: Windows (x64/ARM64), Linux (x64/ARM64), macOS (Apple Silicon / Intel).
- **Prerequisites for running binary**: None (no Node.js, Electron, or CGO runtime required).

---

## 2. Installing from Pre-built Binary

### Windows
1. Download the latest `dbmw-windows-amd64.exe` from GitHub Releases.
2. Rename the binary to `dbmw.exe`.
3. Move `dbmw.exe` to a directory in your system `%PATH%` (e.g. `C:\tools\bin` or `C:\Windows\System32`).

### macOS & Linux
1. Download the release binary matching your CPU architecture (e.g. `dbmw-linux-amd64` or `dbmw-darwin-arm64`).
2. Make the binary executable:
   ```bash
   chmod +x dbmw-linux-amd64
   ```
3. Move it to `/usr/local/bin`:
   ```bash
   sudo mv dbmw-linux-amd64 /usr/local/bin/dbmw
   ```

---

## 3. Installing / Building from Source

If you prefer compiling directly from the repository source code:

### Prerequisites:
- **Go 1.24+**
- **Node.js 20+** & **npm 10+** (used during compilation to build the frontend SPA)

### Compilation Steps:

```bash
# 1. Clone the repository
git clone https://github.com/m-code/dbmw.git
cd dbmw

# 2. Build the embedded frontend assets
cd frontend
npm install
npm run build
cd ..

# 3. Compile the standalone Go binary
go build -o dbmw.exe .
```

---

## 4. Verifying Installation

Verify that the CLI is ready and check system health:

```bash
# Check version information
./dbmw.exe version

# Run the automated diagnostic health check
./dbmw.exe doctor
```

The output will confirm file storage permissions (`~/.dbmw/`), config files, port readiness, and database driver registrations.
