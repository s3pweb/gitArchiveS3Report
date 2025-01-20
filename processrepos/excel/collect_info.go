package excel

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitUtils "github.com/s3pweb/gitArchiveS3Report/utils/git"
	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
	"github.com/s3pweb/gitArchiveS3Report/utils/structs"

	"github.com/alitto/pond"
)

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

	replacements, err := ReadNameReplacement(".config")
	if err != nil {
		logger.Error("Error reading .config file: %v", err)
		return nil, err
	}

	config, err := ReadConfig(".config")
	if err != nil {
		logger.Error("Error reading .config file: %v", err)
		return nil, err
	}

	filesToSearch := config["FILES_TO_SEARCH"]
	termsToSearch := config["TERMS_TO_SEARCH"]

	for _, branchName := range branches {
		logger.Info("Processing branch: %s in repository: %s", branchName, path)

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

		logger.Success("Checked out branch: %s in repository done : %s", branchName, path)

		if err != nil {
			logger.Error("Failed to checkout branch: %s in repository: %s [%s]", branchName, path, err)
			return nil, err
		}

		// Variables pour stocker les informations du commit
		var lastDeveloper string
		var lastCommitDate time.Time
		var commitNbr int
		var lastDeveloperPercentage float64
		var topDeveloper string
		var topDeveloperPercentage float64

		if isShallow {
			// Pour un shallow clone, on prend uniquement le dernier commit
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
			// Utiliser les fonctions existantes pour un clone complet
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
		for _, file := range filesToSearch {
			filesToSearchMap[file] = fileExistsIgnoreCase(path, file)
		}

		termsToSearchMap := make(map[string]bool)
		for _, term := range termsToSearch {
			termsToSearchMap[term] = searchInFiles(path, term)
		}

		trueCountFiles := countTrueInMap(filesToSearchMap)
		trueCountTerms := countTrueInMap(termsToSearchMap)
		totalSearchItems := len(filesToSearchMap) + len(termsToSearchMap)
		trueCount := trueCountFiles + trueCountTerms
		count := fmt.Sprintf("%d/%d", trueCount, totalSearchItems)

		logger.Success("Checked developers info in branch: %s in repository done: %s", branchName, path)

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
			Count:                   count,
			IsShallow:               isShallow,
			CloneDepth:              cloneDepth,
		})
	}
	return infos, nil
}

func CollectBranchInfo(basePath string, logger *logger.Logger) ([]structs.BranchInfo, error) {
	var branchesInfo []structs.BranchInfo

	nbThreads := GetCPU(".config") //runtime.NumCPU() / 2

	logger.Info("Number of threads: %d", nbThreads)
	//time.Sleep(2 * time.Second)

	var mutex sync.Mutex
	pool := pond.New(nbThreads, 0, pond.MinWorkers(nbThreads))

	folders, err := os.ReadDir(basePath)

	if err != nil {
		return nil, err
	}

	for _, oneFolder := range folders {
		path := basePath + "/" + oneFolder.Name()
		logger.Trace("Processing entry: %s", path)

		if oneFolder.IsDir() && isGitRepo(path) {

			pool.Submit(func() {
				infos, err := CollectBranchInfoForOneRepo(logger, branchesInfo, path)

				if err != nil {
					//logger.Error("Error processing repository: %s [%s], continue ...", path, err)
					return
				}

				logger.Success("Processed repository: %s", path)
				logger.Info("Branch info: %v", infos)

				mutex.Lock()
				branchesInfo = append(branchesInfo, infos...)
				mutex.Unlock()
			})
		} else {
			logger.Trace("Not a git repository: %s", path)
		}
	}
	pool.StopAndWait()

	sort.Slice(branchesInfo, func(i, j int) bool {
		if branchesInfo[i].RepoName == branchesInfo[j].RepoName {
			return branchesInfo[i].LastCommitDate.After(branchesInfo[j].LastCommitDate)
		}
		return branchesInfo[i].RepoName < branchesInfo[j].RepoName
	})

	return branchesInfo, nil
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
	if days < 7 {
		return fmt.Sprintf("%d days", days)
	} else if days < 30 {
		weeks := days / 7
		remainingDays := days % 7
		if remainingDays == 0 {
			return fmt.Sprintf("%d weeks", weeks)
		}
		return fmt.Sprintf("%d weeks and %d days", weeks, remainingDays)
	} else {
		months := days / 30
		remainingDays := days % 30
		weeks := remainingDays / 7
		remainingDays = remainingDays % 7
		if weeks == 0 && remainingDays == 0 {
			return fmt.Sprintf("%d months", months)
		} else if weeks == 0 {
			return fmt.Sprintf("%d months and %d days", months, remainingDays)
		} else if remainingDays == 0 {
			return fmt.Sprintf("%d months and %d weeks", months, weeks)
		}
		return fmt.Sprintf("%d months, %d weeks and %d days", months, weeks, remainingDays)
	}
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

// ReadNameReplacement reads the .config file and extracts name replacements from the DEVELOPERS_MAP line
func ReadNameReplacement(filePath string) (map[string]string, error) {
	replacements := make(map[string]string)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var developersMapLine string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "DEVELOPERS_MAP=") {
			developersMapLine = strings.TrimPrefix(line, "DEVELOPERS_MAP=")
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if developersMapLine == "" {
		return nil, nil // No DEVELOPERS_MAP line found
	}

	// Parse the name replacements
	replacementsList := strings.Split(developersMapLine, ";")
	for _, replacement := range replacementsList {
		parts := strings.Split(replacement, "=")
		if len(parts) == 2 {
			replacements[strings.TrimSpace(parts[1])] = strings.TrimSpace(parts[0])
		}
	}

	return replacements, nil
}

// GetCPU reads the .config file and extracts the number of threads to use
func GetCPU(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		return 1
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "CPU=") {
			cpu, err := strconv.Atoi(strings.TrimPrefix(line, "CPU="))
			if err != nil {
				return 1
			}
			return cpu
		}
	}

	if err := scanner.Err(); err != nil {
		return 1
	}

	return 1 // No CPU line found
}
