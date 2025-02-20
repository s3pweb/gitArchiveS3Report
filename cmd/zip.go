package cmd

import (
	"fmt"
	"os"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/processrepos"
	"github.com/spf13/cobra"
)

var (
	zipSourcePath  string
	zipDestPath    string
	uploadToS3     bool
	deleteAfterZip bool
)

var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "Zip a specified path and optionally upload to S3",
	Long: `Zip a specified path (directory or file).
			Creates a single zip file for the specified path with timestamp.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		if zipDestPath == "" {
			zipDestPath = cfg.App.DestDir
			if zipDestPath == "" {
				zipDestPath = "./archive"
			}
		}
		if zipSourcePath == "" {
			return fmt.Errorf("please specify a source path using the --dir-path flag")
		}

		// Check if --remove is used without --upload
		if deleteAfterZip && !uploadToS3 {
			return fmt.Errorf("the --remove option requires --upload to be specified")
		}

		// Step 1: Create the zip file
		fmt.Printf("Creating zip archive from: %s\nDestination: %s\n", zipSourcePath, zipDestPath)
		err := processrepos.Onlyzip(zipSourcePath, zipDestPath)
		if err != nil {
			return fmt.Errorf("error creating zip archive: %v", err)
		}

		// If upload option is not specified, we're done
		if !uploadToS3 {
			return nil
		}

		// Step 2: Find the created zip file
		zipFiles, err := findZipFiles(zipDestPath)
		if err != nil {
			return fmt.Errorf("error finding zip file: %v", err)
		}

		if len(zipFiles) == 0 {
			return fmt.Errorf("no zip file was created in %s", zipDestPath)
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
		if deleteAfterZip {
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
	zipCmd.Flags().StringVarP(&zipSourcePath, "dir-path", "p", "", "Source path to zip (directory or file) (required)")
	zipCmd.Flags().StringVarP(&zipDestPath, "dest-path", "d", "", "Destination directory to save the zip file (default: DEST_DIR in .env if set, otherwise ./archive)")
	zipCmd.Flags().BoolVarP(&uploadToS3, "upload", "u", false, "Upload the created zip file to S3 (requires AWS credentials) (default: false)")
	zipCmd.Flags().BoolVarP(&deleteAfterZip, "remove", "r", false, "Delete local zip file after successful upload (requires --upload) (default: false)")
	zipCmd.MarkFlagRequired("dir-path")
	rootCmd.AddCommand(zipCmd)
}
