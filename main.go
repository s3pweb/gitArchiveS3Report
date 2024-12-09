package main

import (
	"github.com/s3pweb/gitArchiveS3Report/cmd"
	"github.com/s3pweb/gitArchiveS3Report/utils/logger"
)

func main() {

	logger, err := logger.NewLogger("main", "trace")

	if err != nil {
		panic(err)
	}

	logger.Info("Starting Backup-Cobra")

	cmd.Execute()
}
