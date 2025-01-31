package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
	"github.com/s3pweb/gitArchiveS3Report/utils/structs"
)

// ReportExcel generates an Excel report for Bitbucket repositories
// Parameters:
//   - basePath: Base directory containing the repositories
//   - cfg: Configuration object containing report settings
//
// Returns:
//   - error: Any error encountered during report generation
func ReportExcel(basePath, dirDest string, devSheets bool) error {
	logger, err := logger.NewLogger("ReportExcel", "info")
	if err != nil {
		return err
	}

	startTime := time.Now()
	logger.Info("Starting Excel report generation...")

	// Count total repositories before processing
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", basePath, err)
	}

	totalRepos := 0
	for _, entry := range entries {
		if entry.IsDir() && isGitRepo(filepath.Join(basePath, entry.Name())) {
			totalRepos++
		}
	}

	logger.Info("Found %d total repositories to analyze", totalRepos)

	// Collect branch information with progress tracking
	branchesInfo, processedRepos, err := CollectBranchInfo(basePath, logger, totalRepos)
	if err != nil {
		if processedRepos < totalRepos {
			logger.Warn("Processed %d/%d repositories before encountering error", processedRepos, totalRepos)
			logger.Info("Continuing report generation with partial data...")
		}
		// Only return error if no repositories were processed
		if processedRepos == 0 {
			return fmt.Errorf("failed to collect branch information: %v", err)
		}
	}

	// Count unique processed repositories
	repoMap := make(map[string]bool)
	for _, info := range branchesInfo {
		repoMap[info.RepoName] = true
	}
	processedReposCount := len(repoMap)

	excelFile, err := CreateExcelFile(branchesInfo)
	if err != nil {
		return fmt.Errorf("failed to create Excel file: %v", err)
	}

	var mainBranches, developBranches []structs.BranchInfo
	for _, branch := range branchesInfo {
		if branch.BranchName == "main" || branch.BranchName == "origin/main" ||
			branch.BranchName == "master" || branch.BranchName == "origin/master" {
			mainBranches = append(mainBranches, branch)
		} else if branch.BranchName == "develop" || branch.BranchName == "origin/develop" {
			developBranches = append(developBranches, branch)
		}
	}

	err = WriteBranchInfoToExcel(excelFile, branchesInfo, mainBranches, developBranches, devSheets)
	if err != nil {
		return fmt.Errorf("failed to write branch info to Excel: %v", err)
	}

	err = SaveExcelFile(excelFile, basePath, logger)
	if err != nil {
		return fmt.Errorf("failed to save Excel file: %v", err)
	}

	duration := time.Since(startTime).Round(time.Second)
	logger.Info("Report generation completed in %s", duration)
	logger.Info("Total repositories found: %d", totalRepos)
	logger.Info("Repositories successfully processed: %d/%d", processedReposCount, totalRepos)
	logger.Info("Total branches analyzed: %d", len(branchesInfo))
	return nil
}
