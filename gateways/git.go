package gateways

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/int128/gradleupdate/infrastructure"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain"
	"github.com/pkg/errors"
)

type Branch struct{}

func (r *Branch) Get(ctx context.Context, b domain.BranchIdentifier) (domain.Branch, error) {
	client := infrastructure.GitHubClient(ctx)
	payload, resp, err := client.Git.GetRef(ctx, b.Repository.Owner, b.Repository.Name, "refs/heads/"+b.Name)
	if resp != nil && resp.StatusCode == 404 {
		return domain.Branch{}, domain.NotFoundError{Cause: err}
	}
	if err != nil {
		return domain.Branch{}, errors.Wrapf(err, "GitHub API returned error")
	}
	return domain.Branch{
		BranchIdentifier: domain.BranchIdentifier{
			Repository: b.Repository,
			Name:       strings.TrimLeft(payload.GetRef(), "refs/heads/"),
		},
		Commit: domain.CommitIdentifier{
			Repository: b.Repository,
			SHA:        payload.GetObject().GetSHA(),
		},
	}, nil
}

func (r *Branch) Create(ctx context.Context, b domain.Branch) (domain.Branch, error) {
	client := infrastructure.GitHubClient(ctx)
	payload, _, err := client.Git.CreateRef(ctx, b.Repository.Owner, b.Repository.Name, &github.Reference{
		Ref:    github.String("refs/heads/" + b.Name),
		Object: &github.GitObject{SHA: github.String(b.Commit.SHA)},
	})
	if err != nil {
		return domain.Branch{}, errors.Wrapf(err, "GitHub API returned error")
	}
	return domain.Branch{
		BranchIdentifier: domain.BranchIdentifier{
			Repository: b.Repository,
			Name:       strings.TrimLeft(payload.GetRef(), "refs/heads/"),
		},
		Commit: domain.CommitIdentifier{
			Repository: b.Repository,
			SHA:        payload.GetObject().GetSHA(),
		},
	}, nil
}

func (r *Branch) UpdateForce(ctx context.Context, b domain.Branch) (domain.Branch, error) {
	client := infrastructure.GitHubClient(ctx)
	payload, _, err := client.Git.UpdateRef(ctx, b.Repository.Owner, b.Repository.Name, &github.Reference{
		Ref:    github.String("refs/heads/" + b.Name),
		Object: &github.GitObject{SHA: github.String(b.Commit.SHA)},
	}, true)
	if err != nil {
		return domain.Branch{}, errors.Wrapf(err, "GitHub API returned error")
	}
	return domain.Branch{
		BranchIdentifier: domain.BranchIdentifier{
			Repository: b.Repository,
			Name:       strings.TrimLeft(payload.GetRef(), "refs/heads/"),
		},
		Commit: domain.CommitIdentifier{
			Repository: b.Repository,
			SHA:        payload.GetObject().GetSHA(),
		},
	}, nil
}

type Commit struct{}

func (r *Commit) Get(ctx context.Context, c domain.CommitIdentifier) (domain.Commit, error) {
	client := infrastructure.GitHubClient(ctx)
	payload, resp, err := client.Git.GetCommit(ctx, c.Repository.Owner, c.Repository.Name, c.SHA)
	if resp != nil && resp.StatusCode == 404 {
		return domain.Commit{}, domain.NotFoundError{Cause: err}
	}
	if err != nil {
		return domain.Commit{}, errors.Wrapf(err, "GitHub API returned error")
	}
	parents := make([]domain.CommitIdentifier, len(payload.Parents))
	for i, parent := range payload.Parents {
		parents[i] = domain.CommitIdentifier{
			Repository: c.Repository,
			SHA:        parent.GetSHA(),
		}
	}
	return domain.Commit{
		CommitIdentifier: domain.CommitIdentifier{
			Repository: c.Repository,
			SHA:        payload.GetSHA(),
		},
		Message: payload.GetMessage(),
		Parents: parents,
		Tree: domain.TreeIdentifier{
			Repository: c.Repository,
			SHA:        payload.GetTree().GetSHA(),
		},
	}, nil
}

func (r *Commit) Create(ctx context.Context, commit domain.Commit) (domain.Commit, error) {
	client := infrastructure.GitHubClient(ctx)
	ghParents := make([]github.Commit, len(commit.Parents))
	for i, parent := range commit.Parents {
		ghParents[i] = github.Commit{SHA: github.String(parent.SHA)}
	}
	ghCommit, _, err := client.Git.CreateCommit(ctx, commit.Repository.Owner, commit.Repository.Name, &github.Commit{
		Message: github.String(commit.Message),
		Tree:    &github.Tree{SHA: github.String(commit.Tree.SHA)},
		Parents: ghParents,
	})
	if err != nil {
		return domain.Commit{}, errors.Wrapf(err, "GitHub API returned error")
	}
	parents := make([]domain.CommitIdentifier, len(ghCommit.Parents))
	for i, ghParent := range ghCommit.Parents {
		parents[i] = domain.CommitIdentifier{
			Repository: commit.Repository,
			SHA:        ghParent.GetSHA(),
		}
	}
	return domain.Commit{
		CommitIdentifier: domain.CommitIdentifier{
			Repository: commit.Repository,
			SHA:        ghCommit.GetSHA(),
		},
		Message: ghCommit.GetMessage(),
		Parents: parents,
		Tree: domain.TreeIdentifier{
			Repository: commit.Repository,
			SHA:        ghCommit.GetTree().GetSHA(),
		},
	}, nil
}

type Tree struct{}

func (r *Tree) Create(ctx context.Context, id domain.RepositoryIdentifier, base domain.TreeIdentifier, files []domain.File) (domain.TreeIdentifier, error) {
	client := infrastructure.GitHubClient(ctx)
	ghEntries := make([]github.TreeEntry, len(files))
	for i, file := range files {
		content := base64.StdEncoding.EncodeToString(file.Content)
		ghBlob, _, err := client.Git.CreateBlob(ctx, id.Owner, id.Name, &github.Blob{
			Content:  github.String(content),
			Encoding: github.String("base64"),
		})
		if err != nil {
			return domain.TreeIdentifier{}, errors.Wrapf(err, "GitHub API returned error")
		}
		ghEntries[i] = github.TreeEntry{
			Path: github.String(file.Path),
			Mode: github.String(file.Mode),
			SHA:  ghBlob.SHA,
		}
	}
	ghTree, _, err := client.Git.CreateTree(ctx, id.Owner, id.Name, base.SHA, ghEntries)
	if err != nil {
		return domain.TreeIdentifier{}, errors.Wrapf(err, "GitHub API returned error")
	}
	return domain.TreeIdentifier{
		Repository: id,
		SHA:        ghTree.GetSHA(),
	}, nil
}
