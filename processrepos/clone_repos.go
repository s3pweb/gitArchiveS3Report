package processrepos

import (
	"bufio"
	"io"
	"os/exec"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
)

func Clonerepos(dirpath string, cfg *config.Config) error {
	if dirpath == "" {
		dirpath = "./repositories/"
	}

	logger, err := logger.NewLogger("Clonerepos", cfg.Logger.Level)
	if err != nil {
		return err
	}

	// Construction de la commande de base
	args := []string{
		"clone",
		cfg.Bitbucket.Workspace,
		"--scm=bitbucket",
		"--bitbucket-username=" + cfg.Bitbucket.Username,
		"--token=" + cfg.Bitbucket.Token,
		"--path=" + dirpath,
	}

	// Ajout de l'option pour cloner uniquement la branche principale si demandé
	if cfg.App.MainBranchOnly {
		args = append(args, "--branch-filter=main,master")
		args = append(args, "--bare=false")
		logger.Info("Cloning only main/master branches")
	}

	// Ajout de l'option pour le shallow clone
	if cfg.App.ShallowClone {
		args = append(args, "--depth=1")
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
		logger.Error("error waiting for command to finish: %v", err)
		return err
	}

	return nil
}
