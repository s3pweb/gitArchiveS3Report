package processrepos

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
)

// Onlyzip creates a zip archive of the specified directory
// The zip filename includes the source name plus timestamp (YYYYMMDD_HHMM)
func Onlyzip(sourcePath, destPath string) error {
	logger, err := logger.NewLogger("OnlyZip", "trace")
	if err != nil {
		panic(err)
	}

	// Ensure destination directory exists
	err = os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		logger.Error("error creating destination directory: %v", err)
		return err
	}

	// Check if source exists
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		logger.Error("error accessing source path: %v", err)
		return err
	}

	// Get current timestamp for the filename
	timestamp := time.Now().Format("2006-01-02_15h04") // YYYY-MM-DD_HHMM

	// Get a meaningful name for the zip file
	var zipName string
	if sourceInfo.IsDir() {
		// Get the base name of the source directory
		zipName = filepath.Base(sourcePath)
		if zipName == "." || zipName == ".." || zipName == "/" {
			// Use only timestamp if we can't get a meaningful name
			zipName = "archive"
		}
	} else {
		// For a single file, use the filename without extension
		zipName = filepath.Base(sourcePath)
		zipName = zipName[:len(zipName)-len(filepath.Ext(zipName))]
	}

	// Combine name and timestamp
	zipFileName := zipName + "_" + timestamp + ".zip"

	// Create the zip file path
	zipFilePath := filepath.Join(destPath, zipFileName)
	logger.Info("Creating archive: %s", zipFilePath)

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		logger.Error("error creating ZIP file: %v", err)
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Define the base path for relative path calculations
	basePath := sourcePath
	if !sourceInfo.IsDir() {
		// If it's a file, use its directory as base path
		basePath = filepath.Dir(sourcePath)
	}

	// Function to add a file or directory to the zip
	addToZip := func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get path relative to the base path for ZIP entry
		relPath, err := filepath.Rel(basePath, filePath)
		if err != nil {
			return err
		}

		// Skip the root directory itself when zipping a directory
		if sourceInfo.IsDir() && relPath == "." {
			return nil
		}

		// Handle directories
		if info.IsDir() {
			// Use forward slashes for ZIP entries
			_, err = zipWriter.Create(filepath.ToSlash(relPath) + "/")
			return err
		}

		// Handle files
		// If we're zipping a single file and this is that file, use just the filename
		// without directory structure
		zipPath := relPath
		if !sourceInfo.IsDir() && filePath == sourcePath {
			zipPath = filepath.Base(sourcePath)
		}

		// Replace OS-specific path separators with forward slashes for ZIP
		zipPath = filepath.ToSlash(zipPath)

		fileToZip, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		zipEntry, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipEntry, fileToZip)
		return err
	}

	// If source is a directory, walk through it and add all files
	if sourceInfo.IsDir() {
		err = filepath.Walk(sourcePath, addToZip)
		if err != nil {
			logger.Error("error adding files to the ZIP: %v", err)
			return err
		}
	} else {
		// Source is a single file, add just that file
		err = addToZip(sourcePath, sourceInfo, nil)
		if err != nil {
			logger.Error("error adding file to the ZIP: %v", err)
			return err
		}
	}

	logger.Info("Successfully created archive: %s", zipFilePath)
	return nil
}
