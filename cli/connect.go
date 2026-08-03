package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"dbmw/connector"
	"dbmw/core/connection"
	"dbmw/storage"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Interactively manage and test database connections",
	RunE: func(cmd *cobra.Command, args []string) error {
		connStore, err := storage.NewConnectionStore()
		if err != nil {
			return err
		}

		connSvc := connection.NewService(connStore, connector.DefaultFactory)
		reader := bufio.NewReader(os.Stdin)

		conns, err := connSvc.ListConnections()
		if err != nil {
			return err
		}

		fmt.Println("⚡ DBMW Database Connection Manager")
		fmt.Println("------------------------------------")
		if len(conns) > 0 {
			fmt.Println("Saved connections:")
			for i, c := range conns {
				fmt.Printf("  [%d] %s (%s) - %s:%d\n", i+1, c.Name, c.Driver, c.Host, c.Port)
			}
			fmt.Println("  [A] Add new connection profile")
			fmt.Print("\nSelect an option or connection number: ")
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)

			if strings.ToUpper(choice) != "A" {
				idx, err := strconv.Atoi(choice)
				if err == nil && idx >= 1 && idx <= len(conns) {
					target := conns[idx-1]
					fmt.Printf("Testing connection '%s'...\n", target.Name)
					if err := connSvc.TestConnection(context.Background(), target); err != nil {
						fmt.Printf("❌ Connection failed: %v\n", err)
						return nil
					}
					fmt.Println("✅ Connection successful!")
					return nil
				}
			}
		}

		// Add new connection interactively
		fmt.Println("\n--- Add New Connection ---")
		fmt.Print("Connection Name: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		fmt.Print("Driver Type (postgres / mysql / mariadb / sqlite) [postgres]: ")
		drv, _ := reader.ReadString('\n')
		drv = strings.TrimSpace(drv)
		if drv == "" {
			drv = "postgres"
		}

		driverType := connection.DriverType(drv)
		cfg := connection.ConnectionConfig{
			Name:   name,
			Driver: driverType,
		}

		if driverType == connection.DriverSQLite {
			fmt.Print("SQLite File Path [./database.sqlite]: ")
			fPath, _ := reader.ReadString('\n')
			fPath = strings.TrimSpace(fPath)
			if fPath == "" {
				fPath = "./database.sqlite"
			}
			cfg.FilePath = fPath
		} else {
			fmt.Print("Host [127.0.0.1]: ")
			host, _ := reader.ReadString('\n')
			host = strings.TrimSpace(host)
			if host == "" {
				host = "127.0.0.1"
			}
			cfg.Host = host

			defPort := 5432
			if driverType == connection.DriverMySQL || driverType == connection.DriverMariaDB {
				defPort = 3306
			}
			fmt.Printf("Port [%d]: ", defPort)
			pStr, _ := reader.ReadString('\n')
			pStr = strings.TrimSpace(pStr)
			if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
				cfg.Port = p
			} else {
				cfg.Port = defPort
			}

			fmt.Print("Username: ")
			user, _ := reader.ReadString('\n')
			cfg.User = strings.TrimSpace(user)

			fmt.Print("Password: ")
			pass, _ := reader.ReadString('\n')
			cfg.Password = strings.TrimSpace(pass)

			fmt.Print("Database Name: ")
			db, _ := reader.ReadString('\n')
			cfg.Database = strings.TrimSpace(db)
		}

		fmt.Println("\nTesting connection...")
		if err := connSvc.TestConnection(context.Background(), cfg); err != nil {
			fmt.Printf("❌ Connection failed: %v\n", err)
			fmt.Print("Do you still want to save this connection profile? (y/N): ")
			ans, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(ans)) != "y" {
				return nil
			}
		} else {
			fmt.Println("✅ Connection successful!")
		}

		saved, err := connSvc.SaveConnection(cfg)
		if err != nil {
			return fmt.Errorf("failed to save: %w", err)
		}
		fmt.Printf("💾 Saved connection profile with ID: %s\n", saved.ID)
		return nil
	},
}
