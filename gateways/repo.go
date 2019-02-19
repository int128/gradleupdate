package gateways

import (
	"context"
	"encoding/base64"

	"github.com/google/go-github/v24/github"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type RepositoryRepository struct {
	dig.In
	Client *github.Client
}

func (r *RepositoryRepository) Get(ctx context.Context, id git.RepositoryID) (*git.Repository, error) {
	repository, _, err := r.Client.Repositories.Get(ctx, id.Owner, id.Name)
	if err != nil {
		if err, ok := err.(*github.ErrorResponse); ok {
			if err.Response.StatusCode == 404 {
				return nil, errors.Wrapf(&repositoryError{error: err, noSuchEntity: true}, "repository %s not found", id)
			}
		}
		return nil, errors.Wrapf(err, "error from GitHub API")
	}
	return &git.Repository{
		ID: git.RepositoryID{
			Owner: repository.GetOwner().GetLogin(),
			Name:  repository.GetName(),
		},
		Description: repository.GetDescription(),
		AvatarURL:   repository.GetOwner().GetAvatarURL(),
		URL:         repository.GetHTMLURL(),
		DefaultBranch: git.BranchID{
			Repository: git.RepositoryID{
				Owner: repository.GetOwner().GetLogin(),
				Name:  repository.GetName(),
			},
			Name: git.BranchName(repository.GetDefaultBranch()),
		},
	}, nil
}

func (r *RepositoryRepository) GetFileContent(ctx context.Context, id git.RepositoryID, path string) (git.FileContent, error) {
	fc, _, _, err := r.Client.Repositories.GetContents(ctx, id.Owner, id.Name, path, nil)
	if err != nil {
		if err, ok := err.(*github.ErrorResponse); ok {
			if err.Response.StatusCode == 404 {
				return nil, errors.Wrapf(&repositoryError{error: err, noSuchEntity: true}, "file %s not found", path)
			}
		}
		return nil, errors.Wrapf(err, "error from GitHub API")
	}
	if fc == nil {
		return nil, errors.Wrapf(&repositoryError{noSuchEntity: true}, "want a file but got a directory %s", path)
	}
	if fc.GetEncoding() != "base64" {
		return git.FileContent([]byte(*fc.Content)), nil
	}
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(*fc.Content)))
	n, err := base64.StdEncoding.Decode(buf, []byte(*fc.Content))
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode base64")
	}
	return git.FileContent(buf[:n]), nil
}

func (r *RepositoryRepository) GetReadme(ctx context.Context, id git.RepositoryID) (git.FileContent, error) {
	fc, _, err := r.Client.Repositories.GetReadme(ctx, id.Owner, id.Name, nil)
	if err != nil {
		if err, ok := err.(*github.ErrorResponse); ok {
			if err.Response.StatusCode == 404 {
				return nil, errors.Wrapf(&repositoryError{error: err, noSuchEntity: true}, "readme not found")
			}
		}
		return nil, errors.Wrapf(err, "error from GitHub API")
	}
	if fc.GetEncoding() != "base64" {
		return git.FileContent([]byte(*fc.Content)), nil
	}
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(*fc.Content)))
	n, err := base64.StdEncoding.Decode(buf, []byte(*fc.Content))
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode base64")
	}
	return git.FileContent(buf[:n]), nil
}

func (r *RepositoryRepository) Fork(ctx context.Context, id git.RepositoryID) (*git.Repository, error) {
	fork, _, err := r.Client.Repositories.CreateFork(ctx, id.Owner, id.Name, &github.RepositoryCreateForkOptions{})
	if err != nil {
		if _, ok := err.(*github.AcceptedError); ok {
			// Fork in progress
		} else {
			return nil, errors.Wrapf(err, "error from GitHub API")
		}
	}
	return &git.Repository{
		ID: git.RepositoryID{
			Owner: fork.GetOwner().GetLogin(),
			Name:  fork.GetName(),
		},
		Description: fork.GetDescription(),
		AvatarURL:   fork.GetOwner().GetAvatarURL(),
		URL:         fork.GetHTMLURL(),
		DefaultBranch: git.BranchID{
			Repository: git.RepositoryID{
				Owner: fork.GetOwner().GetLogin(),
				Name:  fork.GetName(),
			},
			Name: git.BranchName(fork.GetDefaultBranch()),
		},
	}, nil
}

func (r *RepositoryRepository) GetBranch(ctx context.Context, id git.BranchID) (*git.Branch, error) {
	branch, _, err := r.Client.Repositories.GetBranch(ctx, id.Repository.Owner, id.Repository.Name, id.Name.String())
	if err != nil {
		if err, ok := err.(*github.ErrorResponse); ok {
			if err.Response.StatusCode == 404 {
				return nil, errors.Wrapf(&repositoryError{error: err, noSuchEntity: true}, "branch %s not found", id)
			}
		}
		return nil, errors.Wrapf(err, "error from GitHub API")
	}
	var parents []git.CommitID
	for _, p := range branch.Commit.Parents {
		parents = append(parents, git.CommitID{
			Repository: id.Repository,
			SHA:        git.CommitSHA(p.GetSHA()),
		})
	}
	return &git.Branch{
		ID: id,
		Commit: git.Commit{
			ID: git.CommitID{
				Repository: id.Repository,
				SHA:        git.CommitSHA(branch.Commit.GetSHA()),
			},
			Parents: parents,
			Tree: git.TreeID{
				Repository: id.Repository,
				SHA:        git.TreeSHA(branch.Commit.Commit.Tree.GetSHA()),
			},
		},
	}, nil
}
