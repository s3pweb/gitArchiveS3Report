package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/s3pweb/gitArchiveS3Report/processrepos"
	"github.com/spf13/cobra"
)

var (
	uploadAll  bool
	uploadLast bool
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload repositories to amazon s3",
	Long: `Upload repositories to amazon s3.
			You can specify the directory path or specific zip file to upload.
			
			Available Options:
			--dir-path, -p: Directory path containing zip files
			--all, -a: Upload all zip files in the specified directory
			--last, -l: Upload only the most recent zip file in the directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if dirpath == "" {
			return fmt.Errorf("please specify a path using the --dir-path flag")
		}

		// Check if the path is a directory or a file
		fileInfo, err := os.Stat(dirpath)
		if err != nil {
			return fmt.Errorf("error accessing path %s: %v", dirpath, err)
		}

		// If it's a file, check if it's a zip file
		if !fileInfo.IsDir() {
			if !strings.HasSuffix(strings.ToLower(dirpath), ".zip") {
				return fmt.Errorf("file %s is not a zip file", dirpath)
			}
			fmt.Printf("Uploading file: %s\n", dirpath)
			return processrepos.Upload(dirpath)
		}

		// It's a directory, find zip files
		zipFiles, err := findZipFiles(dirpath)
		if err != nil {
			return fmt.Errorf("error finding zip files: %v", err)
		}

		if len(zipFiles) == 0 {
			return fmt.Errorf("no zip files found in directory %s", dirpath)
		}

		if len(zipFiles) > 1 && !uploadAll && !uploadLast {
			fmt.Printf("Multiple zip files found in %s. Available options:\n", dirpath)
			fmt.Println("  --all, -a: Upload all zip files")
			fmt.Println("  --last, -l: Upload only the most recent zip file")
			fmt.Println("  Or specify the full path to a specific zip file")
			for _, zip := range zipFiles {
				fmt.Printf("  - %s\n", zip)
			}
			return nil
		}

		if uploadLast {
			mostRecent, err := findMostRecentZip(zipFiles)
			if err != nil {
				return fmt.Errorf("error finding most recent zip: %v", err)
			}
			fmt.Printf("Uploading most recent zip file: %s\n", mostRecent)
			return processrepos.Upload(mostRecent)
		}

		if uploadAll {
			fmt.Printf("Uploading all %d zip files from: %s\n", len(zipFiles), dirpath)
			for _, zip := range zipFiles {
				fmt.Printf("Uploading file: %s\n", zip)
				err := processrepos.Upload(zip)
				if err != nil {
					return fmt.Errorf("error uploading file %s: %v", zip, err)
				}
			}
			return nil
		}

		// Default case: single zip in directory
		fmt.Printf("Uploading zip file: %s\n", zipFiles[0])
		return processrepos.Upload(zipFiles[0])
	},
}

// findZipFiles returns all zip files in the given directory
func findZipFiles(dir string) ([]string, error) {
	var zipFiles []string

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".zip") {
			zipFiles = append(zipFiles, filepath.Join(dir, file.Name()))
		}
	}

	return zipFiles, nil
}

// findMostRecentZip returns the most recently modified zip file
func findMostRecentZip(zipFiles []string) (string, error) {
	if len(zipFiles) == 0 {
		return "", fmt.Errorf("no zip files provided")
	}

	var mostRecent string
	var mostRecentTime int64

	for _, zip := range zipFiles {
		fileInfo, err := os.Stat(zip)
		if err != nil {
			return "", err
		}

		modTime := fileInfo.ModTime().Unix()
		if modTime > mostRecentTime || mostRecent == "" {
			mostRecent = zip
			mostRecentTime = modTime
		}
	}

	return mostRecent, nil
}

func init() {
	uploadCmd.Flags().StringVarP(&dirpath, "dir-path", "p", "", "Directory path containing zip files or path to a specific zip file")
	uploadCmd.Flags().BoolVarP(&uploadAll, "all", "a", false, "Upload all zip files in the directory")
	uploadCmd.Flags().BoolVarP(&uploadLast, "last", "l", false, "Upload only the most recent zip file")
	uploadCmd.MarkFlagRequired("dir-path")
	rootCmd.AddCommand(uploadCmd)
}
