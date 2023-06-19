package waltransport

import (
	"fmt"

	"github.com/CAMELNINGA/WAL-transport.git/internal/usecase"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send config",
	Long:  `send config`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res, err := usecase.SendConfig(cmd.Context(), args[0])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
