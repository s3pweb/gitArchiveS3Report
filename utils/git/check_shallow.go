package gitUtils

import (
	"os"
	"path/filepath"
)

// IsShallowClone vérifie si un dépôt est un shallow clone
func IsShallowClone(repoPath string) bool {
	// Dans un shallow clone, Git crée un fichier .git/shallow
	shallowFile := filepath.Join(repoPath, ".git", "shallow")
	_, err := os.Stat(shallowFile)
	return err == nil
}

// GetRepoDepth retourne la profondeur du dépôt (-1 si ce n'est pas un shallow clone)
func GetRepoDepth(repoPath string) int {
	if !IsShallowClone(repoPath) {
		return -1 // -1 indique un clone complet
	}
	return 1 // Pour l'instant, nous supportons uniquement depth=1
}
