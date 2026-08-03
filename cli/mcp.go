package cli

import (
	"fmt"

	"dbmw/connector"
	"dbmw/core/connection"
	"dbmw/core/explorer"
	"dbmw/core/query"
	"dbmw/mcp"
	"dbmw/storage"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start stdio-based read-only Model Context Protocol (MCP) server",
	RunE: func(cmd *cobra.Command, args []string) error {
		connStore, err := storage.NewConnectionStore()
		if err != nil {
			return fmt.Errorf("storage error: %w", err)
		}
		histStore, err := storage.NewHistoryStore()
		if err != nil {
			return fmt.Errorf("history error: %w", err)
		}
		defer histStore.Close()

		connSvc := connection.NewService(connStore, connector.DefaultFactory)
		defer connSvc.CloseAll()

		expSvc := explorer.NewService()
		qSvc := query.NewService(histStore)

		return mcp.RunStdio(connSvc, expSvc, qSvc)
	},
}
