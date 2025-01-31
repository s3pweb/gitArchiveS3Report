package cmd

import (
	"fmt"
	"os"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/processrepos/excel"
	"github.com/spf13/cobra"
)

var (
	devSheets bool
	dirDest   string
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate an Excel report",
	Long: `Generate an Excel report with the following information:
			- Branches
			- Main branches
			- Develop branches
			- Files and terms to search in each branch
			You can specify the directory path where the repositories are cloned (-p, --dir-path).
			Use --dev-sheets or -d to include per-developer sheets in the report.`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		cfg.App.DevSheets = devSheets

		err := excel.ReportExcel(cfg.App.Dir, dirDest, devSheets)
		if err != nil {
			fmt.Printf("Error generating Excel report: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	reportCmd.Flags().StringVarP(&dirDest, "dir-path", "p", "", "Directory where you want the report to be generated (default: value of DIR in .env)")
	reportCmd.Flags().BoolVarP(&devSheets, "dev-sheets", "d", false, "Include developer sheets in the report")
	rootCmd.AddCommand(reportCmd)
}
