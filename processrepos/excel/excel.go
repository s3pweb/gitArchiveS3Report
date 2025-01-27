package excel

import (
	"fmt"
	"time"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
	"github.com/s3pweb/gitArchiveS3Report/utils/structs"
)

// ReportExcel generates an Excel report with the branch information
// and saves it in the specified directory path
// If no directory path is specified, the report will be saved in ./repositories/ (-d, --dir-path)
func ReportExcel(basePath string, cfg *config.Config) error {
	logger, err := logger.NewLogger("ReportExcel", "info")
	if err != nil {
		return err
	}

	if basePath == "" {
		basePath = "./repositories/" + cfg.Bitbucket.Workspace + "/"
	}

	startTime := time.Now()
	logger.Info("Starting Excel report generation...")

	branchesInfo, err := CollectBranchInfo(basePath, logger)
	if err != nil {
		return err
	}

	// Count unique repositories
	repoMap := make(map[string]bool)
	for _, info := range branchesInfo {
		repoMap[info.RepoName] = true
	}
	totalRepos := len(repoMap)

	logger.Info("Analyzing %d repositories...", totalRepos)

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

	err = WriteBranchInfoToExcel(excelFile, branchesInfo, mainBranches, developBranches, cfg.App.DevSheets)
	if err != nil {
		return fmt.Errorf("failed to write branch info to Excel: %v", err)
	}

	err = SaveExcelFile(excelFile, basePath, logger)
	if err != nil {
		return fmt.Errorf("failed to save Excel file: %v", err)
	}

	duration := time.Since(startTime).Round(time.Second)
	logger.Info("Report generation completed in %s", duration)
	logger.Info("Total repositories processed: %d", totalRepos)
	logger.Info("Total branches analyzed: %d", len(branchesInfo))
	logger.Info("Report saved successfully")

	return nil
}
