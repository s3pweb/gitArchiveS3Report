package cmd

import (
	"fmt"
	"os"

	"github.com/s3pweb/gitArchiveS3Report/processrepos"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone repositories from a BitBucket workspace",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		err := processrepos.Clonerepos(dirpath)
		if err != nil {
			fmt.Printf("Error cloning repository: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	cloneCmd.Flags().StringVarP(&dirpath, "dir-path", "d", "", "The directory path where the repositories will be cloned.")
	rootCmd.AddCommand(cloneCmd)
}
