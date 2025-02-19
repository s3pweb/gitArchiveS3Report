package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Git Report Archive S3",
	Short: "Git Report Archive S3: A BitBucket repository backup utility",
	Long: `Git Report Archive S3 is a simple and fast utility for cloning
        and backing up BitBucket repositories from Bitbucket.
        Use it to automate your BitBucket backups with simple commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Définir les styles de couleur
		titleColor := color.New(color.FgHiCyan, color.Bold)
		cmdColor := color.New(color.FgCyan)
		descColor := color.New(color.FgWhite)

		// Afficher le titre
		fmt.Println()
		titleColor.Println("🚀 Welcome to Git Report Archive S3 Tool")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()

		// Afficher les commandes disponibles
		displayCommand(cmdColor, descColor, "clone", "Clone Bitbucket repositories")
		displayCommand(cmdColor, descColor, "report", "Generate Excel report for repositories")
		displayCommand(cmdColor, descColor, "zip", "Create ZIP archives of repositories")
		displayCommand(cmdColor, descColor, "upload", "Upload files to Amazon S3")
		displayCommand(cmdColor, descColor, "zipload", "Create ZIP archive and upload to S3")

		// Afficher l'aide
		fmt.Println()
		descColor.Println("Use './git-archive-s3 [command] --help' for more information about a command.")
		fmt.Println()
	},
}

func displayCommand(cmdColor, descColor *color.Color, cmd, desc string) {
	fmt.Printf("  ")
	cmdColor.Printf("%-10s", cmd)
	descColor.Printf("  %s\n", desc)
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
