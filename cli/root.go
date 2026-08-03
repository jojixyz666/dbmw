package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dbmw/connector"
	"dbmw/core/connection"
	"dbmw/core/data"
	"dbmw/core/erd"
	"dbmw/core/explorer"
	"dbmw/core/project"
	"dbmw/core/query"
	"dbmw/storage"
	"dbmw/web"

	"github.com/spf13/cobra"
)

var (
	portFlag     int
	noOpenFlag   bool
	Version      = "0.0.1"
	CommitHash   = "dev"
	BuildDate    = "unknown"
)

// RootCmd is the primary CLI command entrypoint.
var RootCmd = &cobra.Command{
	Use:   "dbmw",
	Short: "DBMW — Lightweight, self-hosted database management workspace",
	Long: `DBMW is a local-first database management workspace supporting PostgreSQL, 
MySQL, MariaDB, and SQLite with Database Explorer, SQL Editor, Spreadsheet Data Grid, and ERD.`,
	RunE: runServer,
}

func init() {
	RootCmd.PersistentFlags().IntVarP(&portFlag, "port", "p", 0, "Server port (default 8085 or ~/.dbmw/config.json)")
	RootCmd.Flags().BoolVarP(&noOpenFlag, "no-open", "n", false, "Do not auto-open the web browser")

	RootCmd.AddCommand(openCmd)
	RootCmd.AddCommand(connectCmd)
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(doctorCmd)
	RootCmd.AddCommand(mcpCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) error {
	cfgStore, err := storage.NewConfigStore()
	if err != nil {
		return fmt.Errorf("failed to init config store: %w", err)
	}

	appCfg, _ := cfgStore.Get()
	port := appCfg.ServerPort
	if port <= 0 {
		port = 8085
	}
	if portFlag > 0 {
		port = portFlag
	}

	connStore, err := storage.NewConnectionStore()
	if err != nil {
		return fmt.Errorf("failed to init connection store: %w", err)
	}

	histStore, err := storage.NewHistoryStore()
	if err != nil {
		return fmt.Errorf("failed to init history store: %w", err)
	}
	defer histStore.Close()

	connSvc := connection.NewService(connStore, connector.DefaultFactory)
	defer connSvc.CloseAll()

	expSvc := explorer.NewService()
	qSvc := query.NewService(histStore)
	dataSvc := data.NewService()
	erdSvc := erd.NewService()
	projSvc := project.NewService()

	srv, err := web.NewServer(port, connSvc, expSvc, qSvc, dataSvc, erdSvc, projSvc, cfgStore)
	if err != nil {
		return fmt.Errorf("failed to build web server: %w", err)
	}

	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	fmt.Printf("\n🚀 DBMW v%s running at %s\n", Version, url)
	fmt.Println("Press Ctrl+C to stop the server.")

	if !noOpenFlag && appCfg.AutoOpenBrowser {
		go func() {
			time.Sleep(300 * time.Millisecond)
			_ = OpenBrowser(url)
		}()
	}

	// Graceful shutdown channel
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stopChan
		fmt.Println("\nShutting down DBMW gracefully...")
		_ = srv.Shutdown()
		os.Exit(0)
	}()

	return srv.Start()
}
