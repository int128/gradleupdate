package domain

import "fmt"

type PullRequestIdentifier struct {
	Repository        RepositoryIdentifier
	PullRequestNumber int
}

func (p *PullRequestIdentifier) String() string {
	return fmt.Sprintf("%s/pulls#%d", p.Repository.String(), p.PullRequestNumber)
}

type PullRequest struct {
	PullRequestIdentifier
	Base  BranchIdentifier
	Head  BranchIdentifier
	Title string
	Body  string
}
