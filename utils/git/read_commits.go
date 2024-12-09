package gitUtils

import (
	"encoding/json"
	"fmt"
	"log"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func ReadCommits(path string) {
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
			fmt.Printf("----------------------------------------\n")
			fmt.Printf("\nBranch: %s\n", branchName)
			fmt.Printf("----------------------------------------\n")

			// Get commit history for the branch
			printCommitsForBranch(repo, ref)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to iterate over branches: %s", err)
	}
}

func printCommitsForBranch(repo *git.Repository, ref *plumbing.Reference) {
	// Get the commit iterator for the specific branch
	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatalf("Failed to get commits for branch %s: %s", ref.Name().Short(), err)
	}

	var commits []Commit

	// Iterate over the commits
	err = commitIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, Commit{
			Hash:    c.Hash.String(),
			Author:  c.Author.Name,
			Email:   c.Author.Email,
			Date:    c.Author.When,
			Message: c.Message,
		})
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to iterate over commits: %s", err)
	}

	//Convert the commit data to JSON
	commitJSON, err := json.MarshalIndent(commits, "", "  ")
	if err != nil {
		log.Fatalf("Failed to convert commits to JSON: %s", err)
	}

	// Print JSON to standard output
	fmt.Println(string(commitJSON))
}
