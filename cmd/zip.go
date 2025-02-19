package cmd

import (
	"fmt"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/processrepos"
	"github.com/spf13/cobra"
)

var (
	path     string
	destPath string
)

var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "Zip a specified folder",
	Long: `Zip a specified folder.
			You can specify the folder path to zip.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		if destPath == "" {
			destPath = cfg.App.DestDir
		}
		if path == "" {
			return fmt.Errorf("please specify a path using the -path flag")
		}
		return processrepos.Onlyzip(path, destPath)
	},
}

func init() {
	zipCmd.Flags().StringVarP(&path, "dir-path", "p", "", "Folder path to zip")
	zipCmd.Flags().StringVarP(&destPath, "dest-path", "d", "", "Destination path to save the zip file")
	rootCmd.AddCommand(zipCmd)
}
