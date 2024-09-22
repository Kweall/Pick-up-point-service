package cmd

import (
	"fmt"
	"homework/commands"
	"os"

	"github.com/spf13/cobra"
)

var accept_returnCmd = &cobra.Command{
	Use:   "ACCEPT_RETURN [clientID] [orderID]",
	Short: "ACCEPT_RETURN",
	Long:  `Accept return from customer`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running ACCEPT_RETURN command with args:", args)
		if err := commands.AcceptReturn(storage, args); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(accept_returnCmd)
}
