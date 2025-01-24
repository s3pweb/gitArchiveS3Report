package gitUtils

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Branches(repo *git.Repository) ([]string, error) {
	var branches []string

	// Obtenir uniquement les branches locales
	branchRefs, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	err = branchRefs.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref.Name().Short())
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over branches: %w", err)
	}

	return branches, nil
}
