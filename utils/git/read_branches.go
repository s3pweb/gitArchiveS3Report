package gitUtils

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func Branches(repo *git.Repository) ([]string, error) {
	var branches []string

	// Get the reference iterator for all branches (local and remote)
	refs, err := repo.References()
	if err != nil {
		return nil, fmt.Errorf("failed to get references: %w", err)
	}

	// Iterate over the references
	err = refs.ForEach(func(ref *plumbing.Reference) error {

		//color.Red("ref: %s %s %s", ref.Name(), ref.Type(), plumbing.HashReference)
		//color.Yellow("%s %s", ref.Name().IsBranch(), ref.Name().IsRemote())

		// Check if the reference is a branch (local or remote)
		if ref.Type() == plumbing.HashReference {
			if ref.Name().IsRemote() {
				fmt.Printf("Remote Branch: %s\n", ref.Name().Short())

				branches = append(branches, ref.Name().Short())
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over references: %w", err)
	}

	return branches, nil
}
