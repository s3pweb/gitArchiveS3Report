package excel

import (
	"bufio"
	"os"
	"strings"

	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
	"github.com/s3pweb/gitArchiveS3Report/utils/structs"
)

// ReportExcel generates an Excel report with the branch information
// and saves it in the specified directory path
// If no directory path is specified, the report will be saved in ./repositories/ (-d, --dir-path)
func ReportExcel(basePath string) error {
	logger, err := logger.NewLogger("ReportExcel", "trace")
	if err != nil {
		panic(err)
	}
	if basePath == "" {
		workspace, err := ReadWorkspaceName(".secrets")
		if err != nil {
			logger.Error("Error reading workspace name from .secrets file: %v", err)
			return err
		}
		basePath = "./repositories/" + workspace + "/"
	}

	branchesInfo, err := CollectBranchInfo(basePath, logger)
	if err != nil {
		return err
	}

	excelFile, err := CreateExcelFile(branchesInfo)
	if err != nil {
		return err
	}

	var mainBranches, developBranches []structs.BranchInfo
	for _, branch := range branchesInfo {
		if branch.BranchName == "main" || branch.BranchName == "origin/main" || branch.BranchName == "master" || branch.BranchName == "origin/master" {
			mainBranches = append(mainBranches, branch)
		} else if branch.BranchName == "develop" || branch.BranchName == "origin/develop" {
			developBranches = append(developBranches, branch)
		}
	}

	err = WriteBranchInfoToExcel(excelFile, branchesInfo, mainBranches, developBranches)
	if err != nil {
		return err
	}

	err = SaveExcelFile(excelFile, basePath, logger)
	if err != nil {
		return err
	}

	return nil
}

// ReadWorkspaceName reads the .secrets file and extracts the BITBUCKET_WORKSPACE value
func ReadWorkspaceName(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "BITBUCKET_WORKSPACE=") {
			return strings.TrimPrefix(line, "BITBUCKET_WORKSPACE="), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", nil // No BITBUCKET_WORKSPACE line found
}

func ReadConfig(filePath string) (map[string][]string, error) {
	config := make(map[string][]string)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "FILES_TO_SEARCH=") {
			files := strings.TrimPrefix(line, "FILES_TO_SEARCH=")
			config["FILES_TO_SEARCH"] = strings.Split(files, ";")
		} else if strings.HasPrefix(line, "TERMS_TO_SEARCH=") {
			terms := strings.TrimPrefix(line, "TERMS_TO_SEARCH=")
			config["TERMS_TO_SEARCH"] = strings.Split(terms, ";")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}
