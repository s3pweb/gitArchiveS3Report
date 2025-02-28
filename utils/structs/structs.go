package structs

import "time"

// BranchInfo represents information about a branch
type BranchInfo struct {
	RepoName                string
	BranchName              string
	LastCommitDate          time.Time
	TimeSinceLastCommit     string
	Commitnbr               int
	HostLine                string
	LastDeveloper           string
	LastDeveloperPercentage float64
	TopDeveloper            string
	TopDeveloperPercentage  float64
	FilesToSearch           map[string]bool
	TermsToSearch           map[string]bool
	Count                   string
	SelectiveCount          string
	IsShallow               bool
	CloneDepth              int
}
