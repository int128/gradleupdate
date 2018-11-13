package pr

import "fmt"

// Repository represents a GitHub repository.
type Repository struct {
	Owner string
	Repo  string
}

func (r *Repository) String() string {
	return r.Owner + "/" + r.Repo
}

// Branch represents a branch on the repository.
type Branch struct {
	Repository
	Branch string
}

func (b *Branch) String() string {
	return b.Repository.String() + ":" + b.Branch
}

// Commit represents a new commit with files.
type Commit struct {
	Message string
	Files   []File
}

// File represents a file in a commit.
type File struct {
	Path           string
	Mode           string
	EncodedContent string
}

// PullRequest represents a new pull request.
type PullRequest struct {
	Base  Branch
	Head  Branch
	Title string
	Body  string
}

func (p *PullRequest) String() string {
	return fmt.Sprintf("Head[%s]->Base[%s]:%s", p.Head, p.Base, p.Title)
}
