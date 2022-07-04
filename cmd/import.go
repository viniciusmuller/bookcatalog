package cmd

import (
	"fmt"

	"github.com/arcticlimer/bookcatalog/business"
	"github.com/arcticlimer/bookcatalog/importers"
	"github.com/spf13/cobra"
)

var target string

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.PersistentFlags().StringVarP(&target, "target", "t", "required", "folder that documents will be recursively imported from")
}

var importCmd = &cobra.Command{
	Use:   "import <path>",
	Short: "import",
	Long:  `Recursively finds and imports PDF and EPUB files within that directory`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := business.GetDb(databasePath)
		if err != nil {
			return fmt.Errorf("error initializing database: %w", err)
		}
		documentsRepo := business.NewDocumentsRepository(db)

		importerConfig := importers.FsImporterConfig{
			LibraryPath: libraryPath,
			ImagesPath:  imagesPath,
		}

		importer := importers.NewFsImporter(importerConfig, documentsRepo)
		result, err := importer.ImportFiles(target)
		if err != nil {
			return fmt.Errorf("error while importing files: %w", err)
		}

		fmt.Println("Importing finished!")
		for ext, results := range result {
			fmt.Printf("Succesfully imported %d %s files, %d warnings and %d failures\n",
				results.Successes, ext, results.Warnings, results.Failures)
		}

		return nil
	},
}
