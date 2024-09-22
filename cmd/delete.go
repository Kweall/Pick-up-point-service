package cmd

import (
	"fmt"
	"homework/commands"
	"os"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "DELETE [orderID]",
	Short: "DELETE",
	Long:  `Return the order to the courier using the orderID and delete the entry from the file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running DELETE command with args:", args)
		if err := commands.Delete(storage, args); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
