package processrepos

import (
	"bufio"
	"fmt"
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
	logger, err := logger.NewLogger("Clonerepos", "debug")
	if err != nil {
		return err
	}

	if dirpath == "" {
		dirpath = "./repositories/"
	}

	// Construire la commande
	cmd := exec.Command("ghorg", "clone", cfg.Bitbucket.Workspace,
		"--scm=bitbucket",
		"--bitbucket-username="+cfg.Bitbucket.Username,
		"--token="+cfg.Bitbucket.Token,
		"--path="+dirpath)

	// Get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("error getting stdout pipe: %v", err)
		return fmt.Errorf("error getting stdout pipe: %v", err)
	}

	// Get stderr pipe
	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("error getting stderr pipe: %v", err)
		return fmt.Errorf("error getting stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		logger.Error("error starting command: %v", err)
		return fmt.Errorf("error starting command: %v", err)
	}

	// Read command output in real-time
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				logger.Error("error reading stdout: %v", err)
				return
			}
			if err == io.EOF {
				break
			}
			logger.Info("stdout: %s", line)
		}
	}()

	// Read stderr in real-time
	go func() {
		reader := bufio.NewReader(stderr)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				logger.Error("error reading stderr: %v", err)
				return
			}
			if err == io.EOF {
				break
			}
			logger.Error("stderr: %s", line)
		}
	}()

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		logger.Error("error during clone: %v", err)
		return fmt.Errorf("error during clone: %v", err)
	}

	logger.Success("Clone completed successfully")
	return nil
}
