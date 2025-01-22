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
	Long: `Generate an Excel report with the following information:
			- Branches
			- Main branches
			- Develop branches
			- Files and terms to search in each branch
			You can specify the directory path where the repositories are cloned (-d, --dir-path).`,
	Args: cobra.MinimumNArgs(0),
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
