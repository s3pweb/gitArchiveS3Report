package excel

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/s3pweb/gitArchiveS3Report/config"
	gitUtils "github.com/s3pweb/gitArchiveS3Report/utils/git"
	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
	"github.com/s3pweb/gitArchiveS3Report/utils/structs"

	"github.com/alitto/pond"
)

// CollectBranchInfoForOneRepo collects information about branches in a given Git repository.
//
// Parameters:
//   - logger: A pointer to a logger.Logger instance for logging messages.
//   - branchesInfo: A slice of structs.BranchInfo to store information about branches.
//   - path: The file path to the Git repository.
//
// Returns:
//   - A slice of structs.BranchInfo containing information about each branch.
//   - An error if any issues occur during the process.
//
// The function performs the following steps:
//  1. Checks if the repository is a shallow clone and retrieves the clone depth.
//  2. Opens the Git repository located at the specified path.
//  3. Retrieves the list of branches in the repository.
//  4. Reads configuration and name replacement information from a ".config" file.
//  5. Iterates over each branch and checks out the branch.
//  6. Collects information about the last commit, the number of commits, and the top developer for each branch.
//  7. Searches for specified files and terms in the repository.
//  8. Appends the collected information to the branchesInfo slice.
//
// The collected information includes:
//   - Repository name
//   - Branch name
//   - Last commit date
//   - Time since the last commit
//   - Number of commits
//   - Host line from the Docker Compose file
//   - Last developer and their contribution percentage
//   - Top developer and their contribution percentage
//   - Presence of specified files and terms
//   - Count of found items
//   - Whether the repository is a shallow clone
//   - Clone depth
func CollectBranchInfoForOneRepo(logger *logger.Logger, branchesInfo []structs.BranchInfo, path string) ([]structs.BranchInfo, error) {
	var infos []structs.BranchInfo

	isShallow := gitUtils.IsShallowClone(path)
	cloneDepth := gitUtils.GetRepoDepth(path)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	branches, err := gitUtils.Branches(repo)
	if err != nil {
		return nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	logger.Trace("Branches: %v", branches)

	localBranches := make(map[string]bool)

	cfg := config.Get()
	replacements := make(map[string]string)
	if cfg.App.DevelopersMap != "" {
		for _, mapping := range strings.Split(cfg.App.DevelopersMap, ";") {
			parts := strings.Split(mapping, "=")
			if len(parts) == 2 {
				replacements[strings.TrimSpace(parts[1])] = strings.TrimSpace(parts[0])
			}
		}
	}

	for _, branchName := range branches {

		if !strings.HasPrefix(branchName, "origin/") {
			localBranches[branchName] = true
		}

		if strings.HasPrefix(branchName, "origin/") {
			localBranchName := strings.TrimPrefix(branchName, "origin/")
			if localBranches[localBranchName] {
				continue
			}
		}

		if strings.HasPrefix(branchName, "origin/") {
			err = worktree.Checkout(&git.CheckoutOptions{
				Branch: plumbing.NewRemoteReferenceName("origin", strings.TrimPrefix(branchName, "origin/")),
			})
		} else {
			err = worktree.Checkout(&git.CheckoutOptions{
				Branch: plumbing.NewBranchReferenceName(branchName),
			})
		}

		if err != nil {
			logger.Error("Failed to checkout branch: %s in repository: %s [%s]", branchName, path, err)
			return nil, err
		}

		var lastDeveloper string
		var lastCommitDate time.Time
		var commitNbr int
		var lastDeveloperPercentage float64
		var topDeveloper string
		var topDeveloperPercentage float64

		if isShallow {
			head, err := repo.Head()
			if err != nil {
				return nil, err
			}
			commit, err := repo.CommitObject(head.Hash())
			if err != nil {
				return nil, err
			}

			lastDeveloper = commit.Author.Name
			if replacement, ok := replacements[lastDeveloper]; ok {
				lastDeveloper = replacement
			}
			lastCommitDate = commit.Author.When
			commitNbr = 1
			lastDeveloperPercentage = 100
			topDeveloper = lastDeveloper
			topDeveloperPercentage = 100
		} else {
			lastDeveloper, lastCommitDate, err = getLastDeveloperExcludingUser(repo, "bitbucket-pipelines", replacements)
			if err != nil {
				return nil, err
			}
			commitNbr, err = countCommits(repo, "bitbucket-pipelines")
			if err != nil {
				return nil, err
			}
			topDeveloper, topDeveloperPercentage, err = getTopDeveloper(repo, "bitbucket-pipelines", replacements)
			if err != nil {
				return nil, err
			}
			lastDeveloperPercentage = calculateDeveloperPercentage(repo, lastDeveloper)
		}

		dockerComposeName := getDockerComposeFileName(path)
		hostLine := getHostLine(path, dockerComposeName)
		timeSinceLastCommit := formatDuration(time.Since(lastCommitDate))

		filesToSearchMap := make(map[string]bool)
		for _, file := range cfg.App.FilesToSearch {
			filesToSearchMap[file] = fileExistsIgnoreCase(path, file)
		}

		termsToSearchMap := make(map[string]bool)
		for _, term := range cfg.App.TermsToSearch {
			termsToSearchMap[term] = searchInFiles(path, term)
		}

		forbiddenFilesMap := make(map[string]bool)
		for _, file := range cfg.App.ForbiddenFiles {
			forbiddenFilesMap[file] = fileExistsIgnoreCase(path, file)
		}

		trueForbiddenCount := countTrueInMap(forbiddenFilesMap)
		totalForbiddenItems := len(forbiddenFilesMap)
		forbiddenCount := fmt.Sprintf("%d/%d", trueForbiddenCount, totalForbiddenItems)

		trueCountFiles := countTrueInMap(filesToSearchMap)
		trueCountTerms := countTrueInMap(termsToSearchMap)
		totalSearchItems := len(filesToSearchMap) + len(termsToSearchMap)
		trueCount := trueCountFiles + trueCountTerms
		count := fmt.Sprintf("%d/%d", trueCount, totalSearchItems)

		selectiveCountMap := make(map[string]bool)
		for _, item := range cfg.App.TermsFilesToCount {
			if val, exists := filesToSearchMap[item]; exists {
				selectiveCountMap[item] = val
			}
			if val, exists := termsToSearchMap[item]; exists {
				selectiveCountMap[item] = val
			}
		}

		selectiveTrueCount := countTrueInMap(selectiveCountMap)
		selectiveTotalCount := len(selectiveCountMap)
		selectiveCount := fmt.Sprintf("%d/%d", selectiveTrueCount, selectiveTotalCount)

		infos = append(infos, structs.BranchInfo{
			RepoName:                filepath.Base(path),
			BranchName:              branchName,
			LastCommitDate:          lastCommitDate,
			TimeSinceLastCommit:     timeSinceLastCommit,
			Commitnbr:               commitNbr,
			HostLine:                hostLine,
			LastDeveloper:           lastDeveloper,
			LastDeveloperPercentage: lastDeveloperPercentage,
			TopDeveloper:            topDeveloper,
			TopDeveloperPercentage:  topDeveloperPercentage,
			FilesToSearch:           filesToSearchMap,
			TermsToSearch:           termsToSearchMap,
			ForbiddenFiles:          forbiddenFilesMap,
			SelectiveCount:          selectiveCount,
			Count:                   count,
			ForbiddenCount:          forbiddenCount,
			IsShallow:               isShallow,
			CloneDepth:              cloneDepth,
		})
	}
	return infos, nil
}

// CollectBranchInfo collects branch information from git repositories located in the specified base path.
// It uses a thread pool to process multiple repositories concurrently.
//
// Parameters:
//   - basePath: The base directory path where the git repositories are located.
//   - logger: A logger instance for logging information, trace, and errors.
//
// Returns:
//   - A slice of BranchInfo structs containing information about the branches in the repositories.
//   - An error if there is an issue reading the directories or processing the repositories.
//
// collect_info.go
func CollectBranchInfo(basePath string, logger *logger.Logger, totalRepos int) ([]structs.BranchInfo, int, error) {
	startTime := time.Now()
	logWithTime := func(format string, args ...interface{}) {
		elapsed := time.Since(startTime).Round(time.Millisecond)
		if strings.Contains(format, "detected") || strings.Contains(format, "Error") || strings.Contains(format, "empty") {
			logger.Warn("[%s] %s", elapsed, fmt.Sprintf(format, args...))
		} else {
			logger.Info("[%s] %s", elapsed, fmt.Sprintf(format, args...))
		}
	}

	cfg := config.Get()
	var branchesInfo []structs.BranchInfo
	processedRepos := 0
	emptyRepos := make([]string, 0)

	nbThreads := cfg.App.CPU
	if nbThreads <= 0 {
		nbThreads = 1
	}

	logWithTime("Using %d threads for processing", nbThreads)

	var mutex sync.Mutex
	pool := pond.New(nbThreads, 0, pond.MinWorkers(nbThreads))

	folders, err := os.ReadDir(basePath)
	if err != nil {
		return nil, 0, err
	}

	// Create buffered channels with precise sizes
	errorChan := make(chan error, len(folders))
	emptyRepoChan := make(chan string, len(folders))

	for _, oneFolder := range folders {
		path := filepath.Join(basePath, oneFolder.Name())

		if oneFolder.IsDir() && isGitRepo(path) {
			pool.Submit(func() {
				// Check if repository is empty
				isEmpty, err := isEmptyRepository(path)
				if err != nil {
					errorChan <- fmt.Errorf("error checking repository %s: %v", path, err)
					return
				}

				if isEmpty {
					emptyRepoChan <- oneFolder.Name()
					mutex.Lock()
					processedRepos++
					logWithTime("empty repository detected: %s", oneFolder.Name())
					mutex.Unlock()
					return
				}

				infos, err := CollectBranchInfoForOneRepo(logger, branchesInfo, path)

				mutex.Lock()
				if err != nil {
					logWithTime("Error processing repository %s: %v", path, err)
					errorChan <- fmt.Errorf("error in repo %s: %v", path, err)
				} else {
					branchesInfo = append(branchesInfo, infos...)
					processedRepos++
					if processedRepos%(totalRepos/10) == 0 || processedRepos == totalRepos {
						logWithTime("Progress: %d/%d repositories processed (%.1f%%)",
							processedRepos, totalRepos,
							float64(processedRepos)/float64(totalRepos)*100)
					}
				}
				mutex.Unlock()
			})
		}
	}

	logWithTime("Waiting for the last repositories to complete processing...")

	// Create a ticker to periodically log progress
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	go func() {
		pool.StopAndWait()
		done <- true
	}()

	// Wait for either completion or ticker
waitLoop:
	for {
		select {
		case <-done:
			ticker.Stop()
			break waitLoop
		case <-ticker.C:
			logWithTime("Still processing final repositories...")
		}
	}

	logWithTime("All repositories processing completed")

	logWithTime("Starting post-processing phase...")

	close(emptyRepoChan)
	for repoName := range emptyRepoChan {
		emptyRepos = append(emptyRepos, repoName)
	}

	if len(emptyRepos) > 0 {
		sort.Strings(emptyRepos)
		logWithTime("Found %d empty repositories:", len(emptyRepos))
		for _, repoName := range emptyRepos {
			logger.Warn("Empty repository: %s", repoName)
		}
	}

	close(errorChan)
	var errors []error
	for err := range errorChan {
		errors = append(errors, err)
	}

	sort.Slice(branchesInfo, func(i, j int) bool {
		if branchesInfo[i].RepoName == branchesInfo[j].RepoName {
			return branchesInfo[i].LastCommitDate.After(branchesInfo[j].LastCommitDate)
		}
		return branchesInfo[i].RepoName < branchesInfo[j].RepoName
	})

	if len(errors) > 0 {
		logWithTime("%d repositories had errors during processing", len(errors))
		for _, err := range errors {
			logger.Warn("Repository processing error: %v", err)
		}
	}

	return branchesInfo, processedRepos, nil
}

// isEmptyRepository checks if a Git repository is empty (no commits)
func isEmptyRepository(path string) (bool, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, err
	}

	// Try to get HEAD reference
	_, err = repo.Head()
	if err != nil {
		if err == plumbing.ErrReferenceNotFound {
			return true, nil // Repository exists but has no commits
		}
		return false, err
	}

	return false, nil
}

// getLastDeveloperExcludingUser finds the last developer excluding a specific user and returns the developer's name and the commit date
func getLastDeveloperExcludingUser(repo *git.Repository, excludeUser string, replacements map[string]string) (string, time.Time, error) {
	commitIter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return "", time.Time{}, err
	}

	var lastDeveloper string
	var lastCommitDate time.Time
	err = commitIter.ForEach(func(c *object.Commit) error {
		if c.Author.Name != excludeUser {
			lastDeveloper = c.Author.Name
			lastCommitDate = c.Committer.When
			return fmt.Errorf("found") // Stop iteration
		}
		return nil
	})
	if err != nil && err.Error() != "found" {
		return "", time.Time{}, err
	}

	if replacement, ok := replacements[lastDeveloper]; ok {
		lastDeveloper = replacement
	}

	return lastDeveloper, lastCommitDate, nil
}

// getTopDeveloper calculates the top developer and their commit percentage, excluding a specific user
func getTopDeveloper(repo *git.Repository, excludeUser string, replacements map[string]string) (string, float64, error) {
	commitIter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return "", 0, err
	}

	developerCount := make(map[string]int)
	totalCommits := 0

	err = commitIter.ForEach(func(c *object.Commit) error {
		if c.Author.Name == excludeUser {
			return nil
		}
		developerCount[c.Author.Name]++
		totalCommits++
		return nil
	})
	if err != nil {
		return "", 0, err
	}

	var topDeveloper string
	var maxCommits int
	for developer, count := range developerCount {
		if count > maxCommits {
			topDeveloper = developer
			maxCommits = count
		}
	}

	if replacement, ok := replacements[topDeveloper]; ok {
		topDeveloper = replacement
	}

	percentage := calculateDeveloperPercentage(repo, topDeveloper)
	return topDeveloper, percentage, nil
}

// calculateDeveloperPercentage calculates the percentage of commits made by a specific developer
func calculateDeveloperPercentage(repo *git.Repository, developer string) float64 {
	commitIter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return 0
	}
	defer commitIter.Close()

	developerCount := 0
	totalCommits := 0

	err = commitIter.ForEach(func(c *object.Commit) error {
		if c.Author.Name == developer {
			developerCount++
		}
		totalCommits++
		return nil
	})
	if err != nil {
		return 0
	}

	if totalCommits == 0 {
		return 0
	}

	percentage := (float64(developerCount) / float64(totalCommits)) * 100
	roundedPercentage := math.Round(percentage*2) / 2

	return roundedPercentage
}

func isGitRepo(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return !os.IsNotExist(err)
}

func fileExistsIgnoreCase(repoPath, fileNameRegex string) bool {
	// Compile the regular expression
	regex, err := regexp.Compile(fileNameRegex) // (?i) makes the regex case-insensitive
	if err != nil {
		fmt.Printf("Error compiling regex: %v\n", err)
		return false
	}

	err = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && regex.MatchString(info.Name()) {
			return fmt.Errorf("found")
		}
		return nil
	})
	return err != nil && strings.HasPrefix(err.Error(), "found")
}

func searchInFiles(repoPath, searchTermRegex string) bool {
	found := false

	// Compile the regular expression
	regex, err := regexp.Compile(searchTermRegex)
	if err != nil {
		fmt.Printf("Error compiling regex: %v\n", err)
		return false
	}

	err = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if regex.MatchString(string(content)) {
				found = true
				return filepath.SkipDir // Stop walking the directory
			}
		}
		return nil
	})
	if err != nil && err != filepath.SkipDir {
		return false
	}
	return found
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)

	// If more than a month
	if days >= 30 {
		months := days / 30
		return fmt.Sprintf("%d months", months)
	}

	// If more than a week
	if days >= 7 {
		weeks := days / 7
		return fmt.Sprintf("%d weeks", weeks)
	}

	// If less than a week
	return fmt.Sprintf("%d days", days)
}

// countCommits  counts the number of commits in a branch excluding those made by a specific user
func countCommits(repo *git.Repository, excludedUser string) (int, error) {

	commitIter, err := repo.Log(&git.LogOptions{})
	totalCommits := 0
	if err != nil {
		return 0, err
	}
	defer commitIter.Close()

	err = commitIter.ForEach(func(c *object.Commit) error {
		if c.Author.Name != excludedUser {
			totalCommits++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return totalCommits, nil
}

// getHostLine reads a file (ignoring case) and extracts all values inside the Host lines
func getHostLine(dirPath, fileName string) string {
	// Find the file ignoring case
	filePath, err := findFileIgnoreCase(dirPath, fileName)
	if err != nil {
		return ""
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()
	// Read the file line by line

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Check if the line contains the Host keyword (ignoring case)
		if strings.Contains(strings.ToLower(line), "host") {
			// Extract all values inside the parentheses
			start := 0
			maxloop := 0
			for {
				start = strings.Index(line, "(")
				end := strings.LastIndex(line, ")")
				if start != -1 && end != -1 && start < end {
					HostString := line[start+1 : end]
					HostString = strings.Replace(HostString, ") Host (", " ", -1)
					return HostString
				}
				if maxloop > 1000 {
					break
				} else {
					maxloop++
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return ""
	}
	return ""
}

// findFileIgnoreCase finds a file in a directory ignoring case
func findFileIgnoreCase(dirPath, fileName string) (string, error) {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.EqualFold(info.Name(), fileName) {
			return fmt.Errorf("found: %s", path)
		}
		return nil
	})

	if err != nil && strings.HasPrefix(err.Error(), "found: ") {
		return strings.TrimPrefix(err.Error(), "found: "), nil
	}

	return "", fmt.Errorf("file not found: %s", fileName)
}

// getDockerComposeFileName retrieves the name of the docker-compose* file in the given directory
func getDockerComposeFileName(dirPath string) string {
	// Use filepath.Glob to find files that match the pattern docker-compose*
	files, err := filepath.Glob(filepath.Join(dirPath, "docker-compose*"))
	if err != nil {
		return "docker-compose.yaml"
	}
	if len(files) == 0 {
		return "docker-compose.yaml"
	}
	// Return the name of the first match
	return filepath.Base(files[0])
}
