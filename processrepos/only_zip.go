package processrepos

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
)

type zipProgress struct {
	totalFiles     int32
	processedFiles int32
}

func Onlyzip(sourceDir, destDir, workspace string) error {
	logger, err := logger.NewLogger("OnlyZip", "info")
	if err != nil {
		return err
	}

	startTime := time.Now()
	logger.Info("Starting ZIP creation process...")
	logger.Info("Source directory: %s", sourceDir)

	// Count total files before starting
	progress := &zipProgress{}
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			atomic.AddInt32(&progress.totalFiles, 1)
		}
		return nil
	})
	if err != nil {
		logger.Error("Error counting files: %v", err)
		return err
	}

	logger.Info("Found %d files to compress", progress.totalFiles)

	// Create destination directory if it doesn't exist
	if destDir == "" {
		destDir = sourceDir
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		logger.Error("Error creating destination directory: %v", err)
		return err
	}

	// Format current time for filename
	currentTime := time.Now()
	zipFileName := fmt.Sprintf("%s_%s_%s.zip",
		workspace,
		currentTime.Format("2006-01-02"),
		currentTime.Format("15-04"))

	// Create the full path for the ZIP file
	zipFilePath := filepath.Join(destDir, zipFileName)
	logger.Info("Creating ZIP file: %s", zipFilePath)

	// Create the ZIP file
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		logger.Error("Error creating ZIP file: %v", err)
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Start a goroutine to log progress
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				processed := atomic.LoadInt32(&progress.processedFiles)
				total := atomic.LoadInt32(&progress.totalFiles)
				if total > 0 {
					percentage := float64(processed) / float64(total) * 100
					logger.Info("Progress: %d/%d files (%.1f%%)", processed, total, percentage)
				}
			}
		}
	}()

	// Read directory entries
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		logger.Error("Error reading directory: %v", err)
		return err
	}

	// Process each entry
	for _, entry := range entries {
		if entry.IsDir() {
			dirPath := filepath.Join(sourceDir, entry.Name())
			err = addDirectoryToZip(zipWriter, dirPath, "", logger, progress)
			if err != nil {
				logger.Error("Error adding directory to ZIP: %v", err)
				return err
			}
			logger.Info("Added directory to ZIP: %s", entry.Name())
		}
	}

	// Stop the progress logging goroutine
	close(done)

	duration := time.Since(startTime).Round(time.Second)
	logger.Info("ZIP creation completed in %s", duration)
	logger.Info("Successfully processed %d files", progress.processedFiles)
	logger.Info("ZIP file created: %s", zipFilePath)
	return nil
}

func addDirectoryToZip(zipWriter *zip.Writer, dirPath, baseInZip string, logger *logger.Logger, progress *zipProgress) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullPath := filepath.Join(dirPath, file.Name())
		if file.IsDir() {
			// Add directory and its contents recursively
			newBase := filepath.Join(baseInZip, file.Name())
			err = addDirectoryToZip(zipWriter, fullPath, newBase, logger, progress)
			if err != nil {
				return err
			}
		} else {
			// Add file to ZIP
			zipPath := filepath.Join(baseInZip, file.Name())
			err = addFileToZip(zipWriter, fullPath, zipPath, progress)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath, zipPath string, progress *zipProgress) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(zipPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	if err == nil {
		atomic.AddInt32(&progress.processedFiles, 1)
	}
	return err
}
