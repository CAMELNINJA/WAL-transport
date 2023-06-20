package waltransport

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "wal-transport",
	Version: version,
	Short:   "wal-transport is a transport data from WAL to Broker Message Queue.",
	Long: `wal-transport is a transport data from WAL to Broker Message Queue.
	And also transport data from Broker Message Queue to Db or other service.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
