package cmd

import (
	"fmt"
	"os"

	"github.com/s3pweb/gitArchiveS3Report/processrepos"
	"github.com/spf13/cobra"
)

var dirpath string

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload repositories to amazon s3",
	Long: `Upload repositories to amazon s3.
			You can specify the directory path where the repositories will be uploaded.`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Uploading file from path: %s\n", dirpath)
		err := processrepos.Upload(dirpath)
		if err != nil {
			fmt.Printf("Error uploading repository: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	uploadCmd.Flags().StringVarP(&dirpath, "dir-path", "p", "", "The directory path you want to uplaod.")
	uploadCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(uploadCmd)
}
