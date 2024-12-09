package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "backup-cobra",
	Short: "Backup-cobra: A BitBucket repository backup utility",
	Long: `Backup-cobra is a simple and fast utility for cloning
        and backing up BitBucket repositories from Bitbucket.
        Use it to automate your BitBucket backups with simple commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		Cyanp := color.New(color.FgCyan).SprintfFunc()
		color.Blue("Welcome to Backup Tool. Use a subcommand to start.")
		fmt.Printf("Use the ")
		fmt.Printf(Cyanp("clone"))
		fmt.Println(" command to clone Bitbucket repositories.")
		fmt.Printf("Use the ")
		fmt.Printf(Cyanp("report"))
		fmt.Println(" command to create a excel report for each repositories.")
		fmt.Printf("Use the ")
		fmt.Printf(Cyanp("zip"))
		fmt.Println(" command to zip the repositories.")
		fmt.Printf("Use the ")
		fmt.Printf(Cyanp("upload"))
		fmt.Println(" command to upload files or directories into your amazon S3 space.")
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
