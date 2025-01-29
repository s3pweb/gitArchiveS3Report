package cmd

import (
	"fmt"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/processrepos"
	"github.com/spf13/cobra"
)

var (
	mainBranchOnly bool
	shallowClone   bool
	dirpath        string
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone repositories from a BitBucket workspace",
	Long: `Clone repositories from a BitBucket workspace.
			You can specify options to:
			- Clone only main/master branches (-m, --main-only)
			- Perform a shallow clone with only the latest commit (-s, --shallow)
			- Specify the directory path where the repositories will be cloned (-d, --dir-path),
			This is useful when you only need to check the current state of files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()

		if dirpath == "" {
			dirpath = cfg.App.DefaultCloneDir
		}

		cfg.App.MainBranchOnly = mainBranchOnly
		cfg.App.ShallowClone = shallowClone

		if shallowClone {
			cmd.Printf("Warning: Shallow clone will limit the ability to analyze commit history and developer statistics.\n")
		}

		err := processrepos.CloneRepos(dirpath, cfg)
		if err != nil {
			return fmt.Errorf("error cloning repository: %v", err)
		}

		return nil
	},
}

func init() {
	cloneCmd.Flags().StringVarP(&dirpath, "dir-path", "p", "", "The directory path where the repositories will be cloned (default: ./repositories)")
	cloneCmd.Flags().BoolVarP(&mainBranchOnly, "main-only", "m", false, "Clone only the default branch (main/master/develop)")
	cloneCmd.Flags().BoolVarP(&shallowClone, "shallow", "s", false, "Perform a shallow clone with only the latest commit")
	rootCmd.AddCommand(cloneCmd)
}
