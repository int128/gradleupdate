package domain

import "strings"

// RepositoryIdentifier represents an identifier to GitHub repository.
type RepositoryIdentifier struct {
	Owner string
	Repo  string
}

func (r *RepositoryIdentifier) String() string {
	return r.Owner + "/" + r.Repo
}

// Repository represents a GitHub repository.
type Repository struct {
	RepositoryIdentifier
	DefaultBranch BranchIdentifier
}

// GitHubRepositoryURL represents URL for a GitHub repository.
type GitHubRepositoryURL string

// ExtractOwnerAndRepo returns owner and repo for the repository.
func (url GitHubRepositoryURL) ExtractOwnerAndRepo() (string, string) {
	s := strings.Split(string(url), "/")
	if len(s) < 2 {
		return "", ""
	}
	return s[len(s)-2], s[len(s)-1]
}

// File represents a file in a commit.
type File struct {
	Path    string
	Mode    string
	Content []byte
}
