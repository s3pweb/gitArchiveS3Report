package processrepos

import (
	"bufio"
	"io"
	"os/exec"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
)

// CloneRepos clones all repositories from the Bitbucket workspace
// into the specified directory path
// If no directory path is specified, the repositories will be cloned into ./repositories/ (-d, --dir-path)
// If the main branch only option is enabled, only the main branch will be cloned (-m, --main-only)
// If the shallow clone option is enabled, a shallow clone will be performed (-s, --shallow)
// with only the latest commit
func CloneRepos(dirpath string, cfg *config.Config) error {
	if dirpath == "" {
		dirpath = "./repositories/"
	}

	logger, err := logger.NewLogger("Clonerepos", cfg.Logger.Level)
	if err != nil {
		return err
	}

	// Build ghorg clone command
	args := []string{
		"clone",
		cfg.Bitbucket.Workspace,
		"--scm=bitbucket",
		"--bitbucket-username=" + cfg.Bitbucket.Username,
		"--token=" + cfg.Bitbucket.Token,
		"--path=" + dirpath,
	}

	// if main branch only option is enabled, add the branch option
	if cfg.App.MainBranchOnly {
		args = append(args, "--branch=main")
	}

	// if shallow clone option is enabled, add the clone-depth option
	if cfg.App.ShallowClone {
		args = append(args, "--clone-depth=1")
		logger.Info("Performing shallow clone (depth=1)")
	}

	cmd := exec.Command("ghorg", args...)

	// Get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("error getting stdout pipe: %v", err)
		return err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		logger.Error("error starting command: %v", err)
		return err
	}

	// Read command output in real-time
	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			logger.Error("error reading command output: %v", err)
			return err
		}
		if err == io.EOF {
			break
		}
		logger.Info("Output: %s", line)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		// Si la commande échoue avec main, essayer avec master
		if cfg.App.MainBranchOnly {
			logger.Info("Failed with 'main' branch, trying 'master'...")
			// Remplacer main par master dans les arguments
			for i, arg := range args {
				if arg == "--branch=main" {
					args[i] = "--branch=master"
					break
				}
			}

			cmd = exec.Command("ghorg", args...)
			stdout, err = cmd.StdoutPipe()
			if err != nil {
				logger.Error("error getting stdout pipe for master: %v", err)
				return err
			}

			if err := cmd.Start(); err != nil {
				logger.Error("error starting command for master: %v", err)
				return err
			}

			reader = bufio.NewReader(stdout)
			for {
				line, err := reader.ReadString('\n')
				if err != nil && err != io.EOF {
					logger.Error("error reading command output for master: %v", err)
					return err
				}
				if err == io.EOF {
					break
				}
				logger.Info("Output: %s", line)
			}

			if err := cmd.Wait(); err != nil {
				logger.Error("error during master clone: %v", err)
				return err
			}
		} else {
			logger.Error("error during clone: %v", err)
			return err
		}
	}

	return nil
}
