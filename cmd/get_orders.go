package cmd

import (
	"fmt"
	"homework/commands"
	"os"

	"github.com/spf13/cobra"
)

var get_ordersCmd = &cobra.Command{
	Use:   "GET_ORDERS [orderID] (limit optionally)",
	Short: "GET_ORDERS",
	Long:  `Get a list of customer orders`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running GET_ORDERS command with args:", args)
		if err := commands.GetOrders(storage, args); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(get_ordersCmd)
}
