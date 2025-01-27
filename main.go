package main

import (
	"github.com/s3pweb/gitArchiveS3Report/cmd"
	"github.com/s3pweb/gitArchiveS3Report/config"
)

func main() {
	// Initialize config but only enforce validation for commands that need it
	config.Init()

	cmd.Execute()
}
