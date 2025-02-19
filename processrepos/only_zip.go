package processrepos

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
)

func Onlyzip(parentDir, destPath string) error {
	logger, err := logger.NewLogger("OnlyZip", "trace")
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		logger.Error("error creating archive directory: %v", err)
		return err
	}
	// Read subfolders in the parent directory
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		logger.Error("error reading directory: %v", err)
		return err
	}

	// Loop through each entry (subfolder or file) in the parent directory
	for _, entry := range entries {
		// Check if the entry is a directory (we only want to zip folders)
		if entry.IsDir() {
			// Create the full path of the subfolder
			dirPath := filepath.Join(parentDir, entry.Name())
			zipFileName := filepath.Join(destPath, entry.Name()+".zip")

			// Create the ZIP file
			zipFile, err := os.Create(zipFileName)
			if err != nil {
				logger.Error("error creating ZIP file: %v", err)
				return err
			}
			defer zipFile.Close()

			zipWriter := zip.NewWriter(zipFile)
			defer zipWriter.Close()

			// Add files and subfolders to the ZIP file
			err = filepath.Walk(dirPath, func(file string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// Add directories as ZIP entries
				if info.IsDir() {
					relativePath, err := filepath.Rel(parentDir, file)
					if err != nil {
						return err
					}

					_, err = zipWriter.Create(filepath.ToSlash(relativePath) + "/")
					if err != nil {
						return err
					}
					return nil
				}

				// Add files to the ZIP file
				fileToZip, err := os.Open(file)
				if err != nil {
					return err
				}
				defer fileToZip.Close()

				relativePath, err := filepath.Rel(parentDir, file)
				if err != nil {
					return err
				}

				zipEntry, err := zipWriter.Create(filepath.ToSlash(relativePath))
				if err != nil {
					return err
				}

				_, err = io.Copy(zipEntry, fileToZip)
				if err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				logger.Error("error adding files to the ZIP: %v", err)
				return err
			}
			logger.Info("Folder zipped: %s\n", zipFileName)
		}
	}

	return nil
}
