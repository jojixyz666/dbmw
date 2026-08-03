package cli

import (
	"fmt"
	"net/http"
	"time"

	"dbmw/storage"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open the UI for a running server or start one",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgStore, err := storage.NewConfigStore()
		if err != nil {
			return err
		}
		appCfg, _ := cfgStore.Get()
		port := appCfg.ServerPort
		if port <= 0 {
			port = 8085
		}
		if portFlag > 0 {
			port = portFlag
		}

		targetURL := fmt.Sprintf("http://127.0.0.1:%d", port)

		// Check if already responding
		client := http.Client{Timeout: 500 * time.Millisecond}
		resp, err := client.Get(targetURL + "/api/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Printf("DBMW server already running at %s\nOpening browser...\n", targetURL)
			return OpenBrowser(targetURL)
		}

		// Otherwise start server
		return runServer(cmd, args)
	},
}
