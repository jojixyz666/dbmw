package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version, git commit, and Go build details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("DBMW — Database Management Workspace\n")
		fmt.Printf("  Version:    v%s\n", Version)
		fmt.Printf("  Commit:     %s\n", CommitHash)
		fmt.Printf("  Build Date: %s\n", BuildDate)
		fmt.Printf("  Go Version: %s\n", runtime.Version())
		fmt.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}
