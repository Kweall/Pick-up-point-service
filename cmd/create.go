package cmd

import (
	"fmt"
	"homework/commands"
	"os"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "CREATE [clientID] [orderID] [date]",
	Short: "CREATE",
	Long:  `Receive the order from the courier and record it in the database`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running CREATE command with args:", args)
		if err := commands.Create(storage, args); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
