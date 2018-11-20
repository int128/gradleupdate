package repositories

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/app/domain"
	"github.com/pkg/errors"
)

type Branch struct {
	GitHub *github.Client
}

func (r *Branch) Get(ctx context.Context, b domain.BranchIdentifier) (domain.Branch, error) {
	payload, resp, err := r.GitHub.Git.GetRef(ctx, b.Owner, b.Repo, "refs/heads/"+b.Branch)
	if resp.StatusCode == 404 {
		return domain.Branch{}, domain.NotFoundError{Cause: err}
	}
	if err != nil {
		return domain.Branch{}, errors.Wrapf(err, "GitHub API returned error")
	}
	return domain.Branch{
		BranchIdentifier: domain.BranchIdentifier{
			RepositoryIdentifier: b.RepositoryIdentifier,
			Branch:               strings.TrimLeft(payload.GetRef(), "refs/heads/"),
		},
		Commit: domain.CommitIdentifier{
			RepositoryIdentifier: b.RepositoryIdentifier,
			SHA:                  payload.GetObject().GetSHA(),
		},
	}, nil
}

func (r *Branch) Create(ctx context.Context, b domain.Branch) (domain.Branch, error) {
	payload, _, err := r.GitHub.Git.CreateRef(ctx, b.Owner, b.Repo, &github.Reference{
		Ref:    github.String("refs/heads/" + b.Branch),
		Object: &github.GitObject{SHA: github.String(b.Commit.SHA)},
	})
	if err != nil {
		return domain.Branch{}, errors.Wrapf(err, "GitHub API returned error")
	}
	return domain.Branch{
		BranchIdentifier: domain.BranchIdentifier{
			RepositoryIdentifier: b.RepositoryIdentifier,
			Branch:               strings.TrimLeft(payload.GetRef(), "refs/heads/"),
		},
		Commit: domain.CommitIdentifier{
			RepositoryIdentifier: b.RepositoryIdentifier,
			SHA:                  payload.GetObject().GetSHA(),
		},
	}, nil
}

func (r *Branch) UpdateForce(ctx context.Context, b domain.Branch) (domain.Branch, error) {
	payload, _, err := r.GitHub.Git.UpdateRef(ctx, b.Owner, b.Repo, &github.Reference{
		Ref:    github.String("refs/heads/" + b.Branch),
		Object: &github.GitObject{SHA: github.String(b.Commit.SHA)},
	}, true)
	if err != nil {
		return domain.Branch{}, errors.Wrapf(err, "GitHub API returned error")
	}
	return domain.Branch{
		BranchIdentifier: domain.BranchIdentifier{
			RepositoryIdentifier: b.RepositoryIdentifier,
			Branch:               strings.TrimLeft(payload.GetRef(), "refs/heads/"),
		},
		Commit: domain.CommitIdentifier{
			RepositoryIdentifier: b.RepositoryIdentifier,
			SHA:                  payload.GetObject().GetSHA(),
		},
	}, nil
}

type Commit struct {
	GitHub *github.Client
}

func (r *Commit) Get(ctx context.Context, c domain.CommitIdentifier) (domain.Commit, error) {
	payload, resp, err := r.GitHub.Git.GetCommit(ctx, c.Owner, c.Repo, c.SHA)
	if resp.StatusCode == 404 {
		return domain.Commit{}, domain.NotFoundError{Cause: err}
	}
	if err != nil {
		return domain.Commit{}, errors.Wrapf(err, "GitHub API returned error")
	}
	parents := make([]domain.CommitIdentifier, len(payload.Parents))
	for i, parent := range payload.Parents {
		parents[i] = domain.CommitIdentifier{
			RepositoryIdentifier: c.RepositoryIdentifier,
			SHA:                  parent.GetSHA(),
		}
	}
	return domain.Commit{
		CommitIdentifier: domain.CommitIdentifier{
			RepositoryIdentifier: c.RepositoryIdentifier,
			SHA:                  payload.GetSHA(),
		},
		Message: payload.GetMessage(),
		Parents: parents,
		Tree: domain.TreeIdentifier{
			RepositoryIdentifier: c.RepositoryIdentifier,
			SHA:                  payload.GetTree().GetSHA(),
		},
	}, nil
}

func (r *Commit) Create(ctx context.Context, base domain.Commit, files []domain.File) (domain.Commit, error) {
	ghEntries := make([]github.TreeEntry, len(files))
	for i, file := range files {
		content := base64.StdEncoding.EncodeToString(file.Content)
		ghBlob, _, err := r.GitHub.Git.CreateBlob(ctx, base.Owner, base.Repo, &github.Blob{
			Content:  github.String(content),
			Encoding: github.String("base64"),
		})
		if err != nil {
			return domain.Commit{}, errors.Wrapf(err, "GitHub API returned error")
		}
		ghEntries[i] = github.TreeEntry{
			Path: github.String(file.Path),
			Mode: github.String(file.Mode),
			SHA:  ghBlob.SHA,
		}
	}

	ghTree, _, err := r.GitHub.Git.CreateTree(ctx, base.Owner, base.Repo, base.SHA, ghEntries)
	if err != nil {
		return domain.Commit{}, errors.Wrapf(err, "GitHub API returned error")
	}

	ghParents := make([]github.Commit, len(base.Parents))
	for i, parent := range base.Parents {
		ghParents[i] = github.Commit{SHA: github.String(parent.SHA)}
	}
	ghCommit, _, err := r.GitHub.Git.CreateCommit(ctx, base.Owner, base.Repo, &github.Commit{
		Message: github.String(base.Message),
		Tree:    &github.Tree{SHA: ghTree.SHA},
		Parents: ghParents,
	})
	if err != nil {
		return domain.Commit{}, errors.Wrapf(err, "GitHub API returned error")
	}
	parents := make([]domain.CommitIdentifier, len(ghCommit.Parents))
	for i, ghParent := range ghCommit.Parents {
		parents[i] = domain.CommitIdentifier{
			RepositoryIdentifier: base.RepositoryIdentifier,
			SHA:                  ghParent.GetSHA(),
		}
	}
	return domain.Commit{
		CommitIdentifier: domain.CommitIdentifier{
			RepositoryIdentifier: base.RepositoryIdentifier,
			SHA:                  ghCommit.GetSHA(),
		},
		Message: ghCommit.GetMessage(),
		Parents: parents,
		Tree: domain.TreeIdentifier{
			RepositoryIdentifier: base.RepositoryIdentifier,
			SHA:                  ghCommit.GetTree().GetSHA(),
		},
	}, nil
}
