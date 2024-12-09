package cmd

import (
	"fmt"

	"github.com/s3pweb/gitArchiveS3Report/processrepos"
	"github.com/spf13/cobra"
)

var path string

var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "Zip a specified folder",
	RunE: func(cmd *cobra.Command, args []string) error {
		if path == "" {
			return fmt.Errorf("please specify a path using the -path flag")
		}
		return processrepos.Onlyzip(path)
	},
}

func init() {
	zipCmd.Flags().StringVarP(&path, "dir-path", "p", "", "Folder path to zip")
	rootCmd.AddCommand(zipCmd)
}
