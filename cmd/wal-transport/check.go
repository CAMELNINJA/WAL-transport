package waltransport

import (
	"fmt"

	"github.com/CAMELNINGA/WAL-transport.git/internal/usecase"
	"github.com/spf13/cobra"
)

var checkConfigCmd = &cobra.Command{
	Use:   "check",
	Short: "check config",
	Long:  `check config`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res, err := usecase.CheckConfig(args[0])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(checkConfigCmd)
}
