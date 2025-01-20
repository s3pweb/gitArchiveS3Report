package main

import (
	"log"

	"github.com/s3pweb/gitArchiveS3Report/cmd"
	"github.com/s3pweb/gitArchiveS3Report/config"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("Erreur d'initialisation de la configuration : %v", err)
	}

	cmd.Execute()
}
