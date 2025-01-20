package cmd

import (
	"fmt"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/processrepos"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone repositories from a BitBucket workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()

		// Si dirpath n'est pas fourni en argument, utiliser la valeur de la config
		if dirpath == "" {
			dirpath = cfg.App.DefaultCloneDir
		}

		err := processrepos.Clonerepos(dirpath, cfg)
		if err != nil {
			return fmt.Errorf("error cloning repository: %v", err)
		}

		return nil
	},
}

var dirpath string

func init() {
	cloneCmd.Flags().StringVarP(&dirpath, "dir-path", "d", "", "The directory path where the repositories will be cloned")
	rootCmd.AddCommand(cloneCmd)
}
