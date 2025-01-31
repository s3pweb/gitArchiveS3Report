package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
			
			By default, the report will be generated at the workspace root level.
			Use -p or --dir-path to specify a different output location.
			Use --dev-sheets or -d to include per-developer sheets in the report.`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		cfg.App.DevSheets = devSheets

		outputDir := cfg.App.Dir
		if dirDest != "" {
			outputDir = dirDest
		}

		workspacePath := filepath.Join(cfg.App.Dir, cfg.Bitbucket.Workspace)

		err := excel.ReportExcel(workspacePath, outputDir, devSheets)
		if err != nil {
			fmt.Printf("Error generating Excel report: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	reportCmd.Flags().StringVarP(&dirDest, "dir-path", "p", "", "Directory where you want the report to be generated (optional)")
	reportCmd.Flags().BoolVarP(&devSheets, "dev-sheets", "d", false, "Include developer sheets in the report")
	rootCmd.AddCommand(reportCmd)
}
