package cmd

import (
	"fmt"
	"os"

	"github.com/s3pweb/gitArchiveS3Report/processrepos/excel"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate an Excel report",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		err := excel.ReportExcel(dirpath)
		if err != nil {
			fmt.Printf("Error generating Excel report: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	reportCmd.Flags().StringVarP(&dirpath, "dir-path", "p", "", "Folder path")
	rootCmd.AddCommand(reportCmd)
}
