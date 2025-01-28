package gitUtils

import (
	"os"
	"path/filepath"
)

// IsShallowClone checks if the repository is a shallow clone
func IsShallowClone(repoPath string) bool {
	// In a shallow clone, the .git/shallow file exists
	shallowFile := filepath.Join(repoPath, ".git", "shallow")
	_, err := os.Stat(shallowFile)
	return err == nil
}

// GetRepoDepth returns the depth of the repository clone
func GetRepoDepth(repoPath string) int {
	if !IsShallowClone(repoPath) {
		return -1 // -1 indicates that the repository is not shallow
	}
	return 1 // Shallow clone depth is always 1
}
