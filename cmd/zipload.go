package cmd

import (
	"fmt"
	"os"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/processrepos"
	"github.com/spf13/cobra"
)

var (
	zipLoadSourcePath string
	zipLoadDestPath   string
	deleteAfterUpload bool
)

var zipLoadCmd = &cobra.Command{
	Use:   "zipload",
	Short: "Zip and upload to S3",
	Long: `Zip a specified path and upload it to Amazon S3 in a single command.
			This combines the functionality of 'zip' and 'upload' commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()

		// Set default destination directory if not specified
		if zipLoadDestPath == "" {
			zipLoadDestPath = cfg.App.DestDir
			if zipLoadDestPath == "" {
				zipLoadDestPath = "./archive"
			}
		}

		if zipLoadSourcePath == "" {
			return fmt.Errorf("please specify a source path using the --src-path flag")
		}

		// Step 1: Zip the source
		fmt.Printf("Creating zip archive from: %s\nDestination: %s\n", zipLoadSourcePath, zipLoadDestPath)
		err := processrepos.Onlyzip(zipLoadSourcePath, zipLoadDestPath)
		if err != nil {
			return fmt.Errorf("error creating zip archive: %v", err)
		}

		// Step 2: Find the created zip file
		// Since we know Onlyzip creates a single file, we can search for the most recent zip
		// in the destination directory
		zipFiles, err := findZipFiles(zipLoadDestPath)
		if err != nil {
			return fmt.Errorf("error finding zip file: %v", err)
		}

		if len(zipFiles) == 0 {
			return fmt.Errorf("no zip file was created in %s", zipLoadDestPath)
		}

		// Find the most recent zip file (which should be the one we just created)
		zipFile, err := findMostRecentZip(zipFiles)
		if err != nil {
			return fmt.Errorf("error finding created zip file: %v", err)
		}

		// Step 3: Upload the zip file
		fmt.Printf("Uploading file: %s\n", zipFile)
		err = processrepos.Upload(zipFile)
		if err != nil {
			return fmt.Errorf("error uploading file: %v", err)
		}

		// Step 4: Delete local file if requested
		if deleteAfterUpload {
			fmt.Println("Deleting local zip file after successful upload...")
			err := os.Remove(zipFile)
			if err != nil {
				fmt.Printf("Warning: Failed to delete %s: %v\n", zipFile, err)
			} else {
				fmt.Printf("Deleted: %s\n", zipFile)
			}
		}

		fmt.Println("Operation completed successfully.")
		return nil
	},
}

func init() {
	zipLoadCmd.Flags().StringVarP(&zipLoadSourcePath, "src-path", "p", "", "Source path to zip (directory or file) (required)")
	zipLoadCmd.Flags().StringVarP(&zipLoadDestPath, "dest-path", "d", "", "Destination directory to save the zip file (optional)")
	zipLoadCmd.Flags().BoolVarP(&deleteAfterUpload, "remove", "r", false, "Delete local zip file after successful upload (default: false) (optional)")
	zipLoadCmd.MarkFlagRequired("src-path")
	rootCmd.AddCommand(zipLoadCmd)
}
