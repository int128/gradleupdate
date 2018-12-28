package domain

import (
	"fmt"
	"strings"
)

// RepositoryIdentifier points to a repository.
type RepositoryIdentifier struct {
	Owner string
	Name  string
}

func (r RepositoryIdentifier) FullName() string {
	return r.Owner + "/" + r.Name
}

func (r RepositoryIdentifier) String() string {
	return r.FullName()
}

// Repository represents a GitHub repository.
type Repository struct {
	RepositoryIdentifier
	DefaultBranch BranchIdentifier
	Description   string
	AvatarURL     string
}

// RepositoryURL represents URL for a GitHub repository.
type RepositoryURL string

// Parse returns owner and repo for the repository.
func (url RepositoryURL) Parse() *RepositoryIdentifier {
	s := strings.Split(string(url), "/")
	if len(s) < 2 {
		return nil
	}
	return &RepositoryIdentifier{s[len(s)-2], s[len(s)-1]}
}

// FileContent represents content of a file.
type FileContent []byte

func (fc FileContent) String() string {
	return string(fc)
}

// File represents a file in a commit.
type File struct {
	Path    string
	Content FileContent
}

// PullRequestIdentifier points to a pull request.
type PullRequestIdentifier struct {
	Repository        RepositoryIdentifier
	PullRequestNumber int
}

func (p PullRequestIdentifier) String() string {
	return fmt.Sprintf("%s/pulls#%d", p.Repository, p.PullRequestNumber)
}

// PullRequest represents a pull request.
type PullRequest struct {
	PullRequestIdentifier
	BaseBranch BranchIdentifier
	HeadBranch BranchIdentifier
	Title      string
	Body       string
}

// BranchIdentifier points to a branch in a repository.
type BranchIdentifier struct {
	Repository RepositoryIdentifier
	Name       string
}

func (b BranchIdentifier) String() string {
	return b.Repository.String() + ":" + b.Name
}

// Name represents a branch in a repository.
type Branch struct {
	BranchIdentifier
	CommitSHA string
}
