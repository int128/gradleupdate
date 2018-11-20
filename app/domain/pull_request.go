package domain

import "fmt"

type PullRequestIdentifier struct {
	RepositoryIdentifier
	PullRequestNumber int
}

func (p *PullRequestIdentifier) String() string {
	return fmt.Sprintf("%s/pulls#%d", p.RepositoryIdentifier, p.PullRequestNumber)
}

type PullRequest struct {
	PullRequestIdentifier
	Base  BranchIdentifier
	Head  BranchIdentifier
	Title string
	Body  string
}
