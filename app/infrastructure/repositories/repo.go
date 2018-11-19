package repositories

import (
	"context"

	"github.com/pkg/errors"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/app/domain"
)

type Repository struct {
	GitHubClient *github.Client
}

func (r *Repository) GetFile(ctx context.Context, id domain.RepositoryIdentifier, path string) (*domain.File, error) {
	fc, _, _, err := r.GitHubClient.Repositories.GetContents(ctx, id.Owner, id.Repo, path, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Got error from GitHub API")
	}
	if fc == nil {
		return nil, errors.Errorf("Expected file but found directory %s", path)
	}
	content, err := fc.GetContent()
	if err != nil {
		return nil, errors.Wrapf(err, "Could not decode content")
	}
	return &domain.File{
		Path:    path,
		Content: content,
	}, nil
}
