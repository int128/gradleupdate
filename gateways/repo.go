package gateways

import (
	"context"
	"encoding/base64"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/infrastructure/interfaces"
	"github.com/pkg/errors"
)

type RepositoryRepository struct {
	GitHubClientFactory infrastructure.GitHubClientFactory
}

func (r *RepositoryRepository) Get(ctx context.Context, id domain.RepositoryID) (*domain.Repository, error) {
	client := r.GitHubClientFactory.New(ctx)
	repository, _, err := client.Repositories.Get(ctx, id.Owner, id.Name)
	if err != nil {
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
	client := r.GitHubClientFactory.New(ctx)
	fc, _, _, err := client.Repositories.GetContents(ctx, id.Owner, id.Name, path, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error from GitHub API")
	}
	if fc == nil {
		return nil, errors.Errorf("wants a file but got a directory %s", path)
	}
	if fc.GetEncoding() != "base64" {
		return domain.FileContent([]byte(*fc.Content)), nil
	}
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(*fc.Content)))
	n, err := base64.StdEncoding.Decode(buf, []byte(*fc.Content))
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode base64")
	}
	content := buf[:n]
	return domain.FileContent(content), nil
}

func (r *RepositoryRepository) Fork(ctx context.Context, id domain.RepositoryID) (*domain.Repository, error) {
	client := r.GitHubClientFactory.New(ctx)
	fork, _, err := client.Repositories.CreateFork(ctx, id.Owner, id.Name, &github.RepositoryCreateForkOptions{})
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
	client := r.GitHubClientFactory.New(ctx)
	branch, _, err := client.Repositories.GetBranch(ctx, id.Repository.Owner, id.Repository.Name, id.Name)
	if err != nil {
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

func (r *RepositoryRepository) IsNotFoundError(err error) bool {
	if resp, ok := errors.Cause(err).(*github.ErrorResponse); ok {
		return resp.Response.StatusCode == 404
	}
	return false
}
