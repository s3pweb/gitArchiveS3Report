package processrepos

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
)

func CloneRepos(dirpath string, cfg *config.Config) error {
	logger, err := logger.NewLogger("CloneRepos", "info")
	if err != nil {
		return err
	}

	if dirpath == "" {
		dirpath = "./repositories/"
	}

	startTime := time.Now()

	// First check if the directory already exists
	args := []string{
		"clone",
		cfg.Bitbucket.Workspace,
		"--scm=bitbucket",
		"--bitbucket-username=" + cfg.Bitbucket.Username,
		"--token=" + cfg.Bitbucket.Token,
		"--path=" + dirpath,
	}

	if cfg.App.ShallowClone {
		args = append(args, "--clone-depth=1")
	}

	// Execute the clone command
	cmd := exec.Command("ghorg", args...)
	err = executeCloneCommand(cmd, logger)
	if err != nil {
		return fmt.Errorf("clone error: %v", err)
	}

	// If MainBranchOnly is set, clean up non-default branches
	if cfg.App.MainBranchOnly {
		err = cleanupNonDefaultBranches(dirpath, logger)
		if err != nil {
			return fmt.Errorf("error cleaning up branches: %v", err)
		}
	}

	duration := time.Since(startTime).Round(time.Second)
	logger.Info("Clone process completed in %s", duration)
	return nil
}

// cleanupNonDefaultBranches deletes all local branches except the default one
func cleanupNonDefaultBranches(basePath string, logger *logger.Logger) error {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		repoPath := filepath.Join(basePath, entry.Name())
		if !isGitRepo(repoPath) {
			continue
		}

		// Obtain the default branch for the repository
		cmd := exec.Command("git", "-C", repoPath, "remote", "show", "origin")
		output, err := cmd.Output()
		if err != nil {
			logger.Error("Failed to get default branch for %s: %v", entry.Name(), err)
			continue
		}

		// Parse the output to get the default branch
		var defaultBranch string
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "HEAD branch:") {
				defaultBranch = strings.TrimSpace(strings.TrimPrefix(line, "  HEAD branch:"))
				break
			}
		}

		if defaultBranch == "" {
			logger.Error("Could not determine default branch for %s", entry.Name())
			continue
		}

		logger.Info("Default branch for %s is %s", entry.Name(), defaultBranch)

		// Checkout la branche par défaut
		cmd = exec.Command("git", "-C", repoPath, "checkout", defaultBranch)
		if err := cmd.Run(); err != nil {
			logger.Error("Failed to checkout default branch for %s: %v", entry.Name(), err)
			continue
		}

		// Supprime toutes les autres branches locales
		cmd = exec.Command("git", "-C", repoPath, "branch", "|", "grep", "-v", defaultBranch, "|", "xargs", "git", "branch", "-D")
		if err := cmd.Run(); err != nil {
			logger.Error("Failed to delete other branches for %s: %v", entry.Name(), err)
			continue
		}

		logger.Info("Successfully cleaned up %s, keeping only %s branch", entry.Name(), defaultBranch)
	}

	return nil
}

// executeCloneCommand handles the execution of the clone command and logging
func executeCloneCommand(cmd *exec.Cmd, logger *logger.Logger) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe error: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe error: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("command start error: %v", err)
	}

	// Log stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			logger.Info("STDOUT: %s", line)
		}
	}()

	// Log stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			// Skip shallow clone warning as it's not relevant for our analysis
			if !strings.Contains(line, "Shallow clone will limit") {
				logger.Info("STDERR: %s", line)
			}
		}
	}()

	return cmd.Wait()
}

func isGitRepo(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return !os.IsNotExist(err)
}
