package main

import (
	"log"

	"github.com/s3pweb/gitArchiveS3Report/cmd"
	"github.com/s3pweb/gitArchiveS3Report/config"
)

func main() {
	// Initialize config but only enforce validation for commands that need it
	if err := config.Init(); err != nil {
		// Log the error but don't exit
		log.Printf("Warning: Configuration error: %v", err)
	}

	cmd.Execute()
}
