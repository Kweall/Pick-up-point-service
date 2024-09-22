package cmd

import (
	"fmt"
	"homework/commands" // Импортируйте ваш пакет с реализацией команды
	"os"

	"github.com/spf13/cobra"
)

var giveCmd = &cobra.Command{
	Use:   "GIVE orderIDs",
	Short: "GIVE",
	Long:  `Issue orders to ONE! client`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running GIVE command with args:", args)
		if err := commands.Give(storage, args); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(giveCmd)
}
