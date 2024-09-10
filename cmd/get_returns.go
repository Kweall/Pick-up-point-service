package cmd

import (
	"fmt"
	"homework/commands" // Импортируйте ваш пакет с реализацией команды
	"os"

	"github.com/spf13/cobra"
)

var get_returnsCmd = &cobra.Command{
	Use:   "GET_RETURNS [orderID]",
	Short: "GET_RETURNS",
	Long:  `Get a list of returns `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running GET_RETURNS command with args:", args)
		if err := commands.GetReturns(storage, args); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(get_returnsCmd)
}
