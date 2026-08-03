package cli

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"dbmw/connector"
	"dbmw/core/connection"
	"dbmw/storage"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Environment and health check (ports, permissions, config integrity)",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🏥 DBMW Doctor — Environment Health Check")
		fmt.Println("-------------------------------------------")

		// 1. Check ~/.dbmw directory
		dbmwDir, err := storage.GetDbmwDir()
		if err != nil {
			fmt.Printf("❌ Failed to resolve/create ~/.dbmw: %v\n", err)
		} else {
			fmt.Printf("✅ Storage directory: %s\n", dbmwDir)
		}

		// 2. Check config & connection file permissions/integrity
		cfgStore, err := storage.NewConfigStore()
		if err != nil {
			fmt.Printf("❌ Config error: %v\n", err)
		} else {
			cfg, err := cfgStore.Get()
			if err != nil {
				fmt.Printf("⚠️  Config file unreadable: %v (defaults applied)\n", err)
			} else {
				fmt.Printf("✅ Config loaded: port=%d, theme=%s\n", cfg.ServerPort, cfg.Theme)
			}
		}

		connStore, err := storage.NewConnectionStore()
		if err != nil {
			fmt.Printf("❌ Connection store error: %v\n", err)
		} else {
			conns, err := connStore.GetAll()
			if err != nil {
				fmt.Printf("⚠️  Connections file error: %v\n", err)
			} else {
				fmt.Printf("✅ Connection profiles: %d saved connection(s)\n", len(conns))
			}
		}

		// 3. Check SQLite history database
		histStore, err := storage.NewHistoryStore()
		if err != nil {
			fmt.Printf("❌ History database error: %v\n", err)
		} else {
			histStore.Close()
			histPath := filepath.Join(dbmwDir, "history.db")
			if fi, err := os.Stat(histPath); err == nil {
				fmt.Printf("✅ History database: %s (%.1f KB)\n", histPath, float64(fi.Size())/1024.0)
			}
		}

		// 4. Check port availability
		port := 8085
		if portFlag > 0 {
			port = portFlag
		}
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			fmt.Printf("⚠️  Port %d is currently occupied or in use\n", port)
		} else {
			ln.Close()
			fmt.Printf("✅ Port %d is available\n", port)
		}

		// 5. Check database connectors registry
		drivers := []connection.DriverType{
			connection.DriverPostgres,
			connection.DriverMySQL,
			connection.DriverMariaDB,
			connection.DriverSQLite,
		}
		for _, d := range drivers {
			_, err := connector.DefaultFactory(d)
			if err != nil {
				fmt.Printf("❌ Driver %s initialization error: %v\n", d, err)
			} else {
				fmt.Printf("✅ Driver supported: %s\n", d)
			}
		}

		fmt.Println("-------------------------------------------")
		fmt.Println("🎉 Diagnostic check complete.")
		return nil
	},
}
