package gitUtils

import (
	"fmt"
	"log"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func ReadFiles(path string) {
	// Open the Git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Fatalf("Failed to open repository: %s", err)
	}

	// Get references (local and remote branches)
	refs, err := repo.References()
	if err != nil {
		log.Fatalf("Failed to get references: %s", err)
	}

	// Iterate over each reference (branch)
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference && (ref.Name().IsBranch() || ref.Name().IsRemote()) {
			branchName := ref.Name().Short()

			// Print branch name
			fmt.Printf("\nBranch: %s\n", branchName)

			// List files in the branch
			listFilesInBranch(repo, ref)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to iterate over branches: %s", err)
	}
}

func listFilesInBranch(repo *git.Repository, ref *plumbing.Reference) {
	// Get the commit object for the branch
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		log.Fatalf("Failed to get commit for branch %s: %s", ref.Name().Short(), err)
	}

	// Get the tree associated with the commit
	tree, err := commit.Tree()
	if err != nil {
		log.Fatalf("Failed to get tree for commit %s: %s", commit.Hash.String(), err)
	}

	// Print files in the tree
	err = tree.Files().ForEach(func(file *object.File) error {
		fmt.Printf("  File: %s %d\n", file.Name, file.Size)

		content, err := file.Contents()
		if err != nil {
			log.Fatalf("Failed to get contents of file %s: %s", file.Name, err)

		}

		// size := len(content)
		// if size < 1000 {
		// 	fmt.Printf("%s\n", content)
		// }

		stringToFind := "sonar"

		hasVault := strings.Contains(
			strings.ToLower(content),
			strings.ToLower(stringToFind),
		)

		if hasVault {
			fmt.Printf("File has %s : %s %d\n", stringToFind, file.Name, file.Size)
			//fmt.Printf("%s\n", content)
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Failed to iterate over files in branch %s: %s", ref.Name().Short(), err)
	}
}
