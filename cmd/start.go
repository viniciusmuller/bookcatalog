package cmd

import (
	"fmt"

	"github.com/arcticlimer/bookcatalog/business"
	"github.com/arcticlimer/bookcatalog/server"
	"github.com/spf13/cobra"
)

var address string

func init() {
	rootCmd.AddCommand(startCommand)
	startCommand.PersistentFlags().StringVarP(&address, "address", "a", ":8080", "address to run the server (e.g: :8080)")
}

var startCommand = &cobra.Command{
	Use:   "start",
	Short: "start",
	Long:  `Starts the server that will import and serve documents`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := business.GetDb(databasePath)
		if err != nil {
			return fmt.Errorf("error initializing database: %w", err)
		}

		return server.Start(db, address, libraryPath, imagesPath)
	},
}
