package processrepos

import (
	"bufio"
	"io"
	"os/exec"

	"github.com/s3pweb/gitArchiveS3Report/utils"
	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
)

func Clonerepos(dirpath string) error {
	if dirpath == "" {
		dirpath = "./repositories/"
	}
	logger, err := logger.NewLogger("Clonerepos", "trace")
	if err != nil {
		return err
	}
	secrets, err := utils.ReadConfigFile(".secrets")
	if err != nil {
		logger.Error("Error reading file: %v", err)
	}
	token := secrets["BITBUCKET_TOKEN"]
	username := secrets["BITBUCKET_USERNAME"]
	workspace := secrets["BITBUCKET_WORKSPACE"]

	cmd := exec.Command("ghorg", "clone", workspace,
		"--scm=bitbucket",
		"--bitbucket-username="+username,
		"--token="+token,
		"--path="+dirpath)

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
		// Read each line as it's available
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			logger.Error("error reading command output: %v", err)
			return err
		}
		// If end of output
		if err == io.EOF {
			break
		}
		// Simply print the line
		logger.Info("Output: %s", line)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		logger.Error("error waiting for command to finish: %v", err)
		return err
	}
	return nil
}
