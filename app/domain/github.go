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

func (r *RepositoryIdentifier) FullName() string {
	return r.Owner + "/" + r.Name
}

func (r *RepositoryIdentifier) String() string {
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

// File represents a file in a commit.
type File struct {
	Path    string
	Mode    string
	Content []byte
}

// PullRequestIdentifier points to a pull request.
type PullRequestIdentifier struct {
	Repository        RepositoryIdentifier
	PullRequestNumber int
}

func (p *PullRequestIdentifier) String() string {
	return fmt.Sprintf("%s/pulls#%d", p.Repository.String(), p.PullRequestNumber)
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

func (b *BranchIdentifier) String() string {
	return b.Repository.String() + ":" + b.Name
}

// Name represents a branch in a repository.
type Branch struct {
	BranchIdentifier
	Commit CommitIdentifier
}

// CommitIdentifier points to a commit in a repository.
type CommitIdentifier struct {
	Repository RepositoryIdentifier
	SHA        string
}

// Commit represents a commit in a repository.
type Commit struct {
	CommitIdentifier
	Message string
	Parents []CommitIdentifier
	Tree    TreeIdentifier
}

func (c *Commit) GetSingleParent() *CommitIdentifier {
	return nil
}

// TreeIdentifier points to a tree in a repository.
type TreeIdentifier struct {
	Repository RepositoryIdentifier
	SHA        string
}
