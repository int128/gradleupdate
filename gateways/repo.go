package gateways

import (
	"context"
	"encoding/base64"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type RepositoryRepository struct {
	dig.In
	Client *github.Client
}

func (r *RepositoryRepository) Get(ctx context.Context, id domain.RepositoryID) (*domain.Repository, error) {
	repository, _, err := r.Client.Repositories.Get(ctx, id.Owner, id.Name)
	if err != nil {
		if err, ok := err.(*github.ErrorResponse); ok {
			if err.Response.StatusCode == 404 {
				return nil, errors.Wrapf(&repositoryError{error: err, noSuchEntity: true}, "repository %s not found", id)
			}
		}
		return nil, errors.Wrapf(err, "error from GitHub API")
	}
	return &domain.Repository{
		ID: domain.RepositoryID{
			Owner: repository.GetOwner().GetLogin(),
			Name:  repository.GetName(),
		},
		Description: repository.GetDescription(),
		AvatarURL:   repository.GetOwner().GetAvatarURL(),
		DefaultBranch: domain.BranchID{
			Repository: domain.RepositoryID{
				Owner: repository.GetOwner().GetLogin(),
				Name:  repository.GetName(),
			},
			Name: repository.GetDefaultBranch(),
		},
	}, nil
}

func (r *RepositoryRepository) GetFileContent(ctx context.Context, id domain.RepositoryID, path string) (domain.FileContent, error) {
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
		return domain.FileContent([]byte(*fc.Content)), nil
	}
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(*fc.Content)))
	n, err := base64.StdEncoding.Decode(buf, []byte(*fc.Content))
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode base64")
	}
	return domain.FileContent(buf[:n]), nil
}

func (r *RepositoryRepository) GetReadme(ctx context.Context, id domain.RepositoryID) (domain.FileContent, error) {
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
		return domain.FileContent([]byte(*fc.Content)), nil
	}
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(*fc.Content)))
	n, err := base64.StdEncoding.Decode(buf, []byte(*fc.Content))
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode base64")
	}
	return domain.FileContent(buf[:n]), nil
}

func (r *RepositoryRepository) Fork(ctx context.Context, id domain.RepositoryID) (*domain.Repository, error) {
	fork, _, err := r.Client.Repositories.CreateFork(ctx, id.Owner, id.Name, &github.RepositoryCreateForkOptions{})
	if err != nil {
		if _, ok := err.(*github.AcceptedError); ok {
			// Fork in progress
		} else {
			return nil, errors.Wrapf(err, "error from GitHub API")
		}
	}
	return &domain.Repository{
		ID: domain.RepositoryID{
			Owner: fork.GetOwner().GetLogin(),
			Name:  fork.GetName(),
		},
		Description: fork.GetDescription(),
		AvatarURL:   fork.GetOwner().GetAvatarURL(),
		DefaultBranch: domain.BranchID{
			Repository: domain.RepositoryID{
				Owner: fork.GetOwner().GetLogin(),
				Name:  fork.GetName(),
			},
			Name: fork.GetDefaultBranch(),
		},
	}, nil
}

func (r *RepositoryRepository) GetBranch(ctx context.Context, id domain.BranchID) (*domain.Branch, error) {
	branch, _, err := r.Client.Repositories.GetBranch(ctx, id.Repository.Owner, id.Repository.Name, id.Name)
	if err != nil {
		if err, ok := err.(*github.ErrorResponse); ok {
			if err.Response.StatusCode == 404 {
				return nil, errors.Wrapf(&repositoryError{error: err, noSuchEntity: true}, "branch %s not found", id)
			}
		}
		return nil, errors.Wrapf(err, "error from GitHub API")
	}
	var parents []domain.CommitID
	for _, p := range branch.Commit.Parents {
		parents = append(parents, domain.CommitID{
			Repository: id.Repository,
			SHA:        domain.CommitSHA(p.GetSHA()),
		})
	}
	return &domain.Branch{
		ID: id,
		Commit: domain.Commit{
			ID: domain.CommitID{
				Repository: id.Repository,
				SHA:        domain.CommitSHA(branch.Commit.GetSHA()),
			},
			Parents: parents,
			Tree: domain.TreeID{
				Repository: id.Repository,
				SHA:        domain.TreeSHA(branch.Commit.Commit.Tree.GetSHA()),
			},
		},
	}, nil
}
