package processrepos

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/s3pweb/gitArchiveS3Report/config"
	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
)

func CloneRepos(dirpath string, cfg *config.Config) error {
	logger, err := logger.NewLogger("Clonerepos", "debug")
	if err != nil {
		return err
	}

	if dirpath == "" {
		dirpath = "./repositories/"
	}

	fmt.Println("Config = ", cfg)

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

	if cfg.App.MainBranchOnly {
		args = append(args, "--branch=master")
	}

	cmd := exec.Command("ghorg", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("erreur pipe stdout: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("erreur pipe stderr: %v", err)
	}

	var errMsgs []string

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("erreur démarrage commande: %v", err)
	}

	go func() {
		reader := bufio.NewReader(stdout)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				return
			}
			if err == io.EOF {
				break
			}
			logger.Info("stdout: %s", strings.TrimSpace(line))
		}
	}()

	go func() {
		reader := bufio.NewReader(stderr)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				return
			}
			if err == io.EOF {
				break
			}
			errLine := strings.TrimSpace(line)
			errMsgs = append(errMsgs, errLine)
			logger.Error("stderr: %s", errLine)
		}
	}()

	if err := cmd.Wait(); err != nil {
		if len(errMsgs) > 0 {
			return fmt.Errorf("erreur clone: %v\nMessages d'erreur:\n%s", err, strings.Join(errMsgs, "\n"))
		}
		return fmt.Errorf("erreur clone: %v", err)
	}

	return nil
}
