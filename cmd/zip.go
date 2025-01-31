package cmd

import (
	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/processrepos"
	"github.com/spf13/cobra"
)

var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "Create a ZIP file with all repositories",
	Long: `Create a ZIP file with all repositories.
			You can specify the directory path where the ZIP file will be created (-p, --dir-path).
			If not specified, the ZIP file will be created in the value of DIR in .env.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()

		return processrepos.Onlyzip(cfg.App.Dir, dirDest, cfg.Bitbucket.Workspace)
	},
}

func init() {
	zipCmd.Flags().StringVarP(&dirDest, "dir-path", "p", "", "The directory path where the ZIP file will be created (default: value of ZIP_DIR in .env)")
	rootCmd.AddCommand(zipCmd)
}
