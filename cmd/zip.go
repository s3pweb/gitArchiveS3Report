package cmd

import (
	"fmt"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/processrepos"
	"github.com/spf13/cobra"
)

var (
	zipSourcePath string
	zipDestPath   string
)

var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "Zip a specified path",
	Long: `Zip a specified path (directory or file).
			Creates a single zip file for the specified path.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		if zipDestPath == "" {
			zipDestPath = cfg.App.DestDir
			if zipDestPath == "" {
				zipDestPath = "./archive"
			}
		}
		if zipSourcePath == "" {
			return fmt.Errorf("please specify a source path using the --src-path (or -p) flag")
		}

		fmt.Printf("Creating zip archive from: %s\nDestination: %s\n", zipSourcePath, zipDestPath)
		return processrepos.Onlyzip(zipSourcePath, zipDestPath)
	},
}

func init() {
	zipCmd.Flags().StringVarP(&zipSourcePath, "src-path", "p", "", "Source path to zip (directory or file) (required)")
	zipCmd.Flags().StringVarP(&zipDestPath, "dest-path", "d", "", "Destination path to save the zip file (default: DEST_DIR if set, otherwise ./archive)")
	rootCmd.AddCommand(zipCmd)
}
