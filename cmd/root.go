// cmd/root.go
package cmd

import (
	"fmt"
	"homework/storage/json_file"
	"os"

	"github.com/spf13/cobra"
)

var (
	storage *json_file.Storage
	rootCmd = &cobra.Command{
		Use:   "console-app",
		Short: "Management of orders and returns.",
		Long:  `A command line application to manage orders and returns for customers.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Инициализируем storage
			var err error
			storage, err = json_file.NewStorage("storage/json_file/data.json")
			if err != nil {
				return fmt.Errorf("can't init storage: %v", err)
			}
			return nil
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
