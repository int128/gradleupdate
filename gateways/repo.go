package gateways

import (
	"context"
	"encoding/base64"

	"github.com/int128/gradleupdate/infrastructure"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain"
	"github.com/pkg/errors"
)

type RepositoryRepository struct {
	GitHubClient *infrastructure.GitHubClientFactory
}

func (r *RepositoryRepository) Get(ctx context.Context, id domain.RepositoryID) (*domain.Repository, error) {
	client := r.GitHubClient.New(ctx)
	repository, resp, err := client.Repositories.Get(ctx, id.Owner, id.Name)
	if resp != nil && resp.StatusCode == 404 {
		return nil, domain.NotFoundError{Cause: err}
	}
	if err != nil {
		return nil, errors.Wrapf(err, "GitHub API returned error")
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
	client := r.GitHubClient.New(ctx)
	fc, _, resp, err := client.Repositories.GetContents(ctx, id.Owner, id.Name, path, nil)
	if resp != nil && resp.StatusCode == 404 {
		return nil, domain.NotFoundError{Cause: err}
	}
	if err != nil {
		return nil, errors.Wrapf(err, "GitHub API returned error")
	}
	if fc == nil {
		return nil, errors.Errorf("Expected file but found directory %s", path)
	}
	if fc.GetEncoding() != "base64" {
		return domain.FileContent([]byte(*fc.Content)), nil
	}
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(*fc.Content)))
	n, err := base64.StdEncoding.Decode(buf, []byte(*fc.Content))
	if err != nil {
		return nil, errors.Wrapf(err, "Could not decode content")
	}
	content := buf[:n]
	return domain.FileContent(content), nil
}

func (r *RepositoryRepository) Fork(ctx context.Context, id domain.RepositoryID) (*domain.Repository, error) {
	client := r.GitHubClient.New(ctx)
	fork, resp, err := client.Repositories.CreateFork(ctx, id.Owner, id.Name, &github.RepositoryCreateForkOptions{})
	if resp != nil && resp.StatusCode == 404 {
		return nil, domain.NotFoundError{Cause: err}
	}
	if err != nil {
		if _, ok := err.(*github.AcceptedError); ok {
			// Fork in progress
		} else {
			return nil, errors.Wrapf(err, "GitHub API returned error")
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
