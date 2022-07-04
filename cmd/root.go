package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	databasePath string
	libraryPath  string
	imagesPath   string
)

// TODO: Create repair commands that does things such as adding dangling files
// to the database
func init() {
	rootCmd.PersistentFlags().StringVarP(&databasePath, "dbpath", "d", "db/bookcatalog.db", "path of the sqlite3 database file")
	rootCmd.PersistentFlags().StringVarP(&libraryPath, "librarypath", "l", "library", "folder where the documents will be stored ")
	rootCmd.PersistentFlags().StringVarP(&imagesPath, "imagespath", "i", "img", "folder where e-book cover images will be stored")
}

var rootCmd = &cobra.Command{
	Use:   "bookcatalog",
	Short: "bookcatalog CLI",
	Long:  `Minimal self-hosted ebook management service`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
