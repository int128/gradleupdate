package repositories

import (
	"context"
	"encoding/base64"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/app/domain"
	"github.com/pkg/errors"
)

type Repository struct {
	GitHub *github.Client
}

func (r *Repository) Get(ctx context.Context, id domain.RepositoryIdentifier) (domain.Repository, error) {
	repository, resp, err := r.GitHub.Repositories.Get(ctx, id.Owner, id.Repo)
	if resp.StatusCode == 404 {
		return domain.Repository{}, domain.NotFoundError{Cause: err}
	}
	if err != nil {
		return domain.Repository{}, errors.Wrapf(err, "GitHub API returned error")
	}
	return domain.Repository{
		RepositoryIdentifier: domain.RepositoryIdentifier{
			Owner: repository.GetOwner().GetLogin(),
			Repo:  repository.GetName(),
		},
		DefaultBranch: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{
				Owner: repository.GetOwner().GetLogin(),
				Repo:  repository.GetName(),
			},
			Branch: repository.GetDefaultBranch(),
		},
	}, nil
}

func (r *Repository) GetFile(ctx context.Context, id domain.RepositoryIdentifier, path string) (domain.File, error) {
	fc, _, resp, err := r.GitHub.Repositories.GetContents(ctx, id.Owner, id.Repo, path, nil)
	if resp.StatusCode == 404 {
		return domain.File{}, domain.NotFoundError{Cause: err}
	}
	if err != nil {
		return domain.File{}, errors.Wrapf(err, "GitHub API returned error")
	}
	if fc == nil {
		return domain.File{}, errors.Errorf("Expected file but found directory %s", path)
	}
	var content []byte
	switch fc.GetEncoding() {
	case "base64":
		buf := make([]byte, base64.StdEncoding.DecodedLen(len(*fc.Content)))
		n, err := base64.StdEncoding.Decode(buf, []byte(*fc.Content))
		if err != nil {
			return domain.File{}, errors.Wrapf(err, "Could not decode content")
		}
		content = buf[:n]
	default:
		content = []byte(*fc.Content)
	}
	return domain.File{
		Path:    path,
		Content: content,
	}, nil
}

func (r *Repository) Fork(ctx context.Context, id domain.RepositoryIdentifier) (domain.Repository, error) {
	fork, resp, err := r.GitHub.Repositories.CreateFork(ctx, id.Owner, id.Repo, &github.RepositoryCreateForkOptions{})
	if resp.StatusCode == 404 {
		return domain.Repository{}, domain.NotFoundError{Cause: err}
	}
	if err != nil {
		if _, ok := err.(*github.AcceptedError); ok {
			// Fork in progress
		} else {
			return domain.Repository{}, errors.Wrapf(err, "GitHub API returned error")
		}
	}
	return domain.Repository{
		RepositoryIdentifier: domain.RepositoryIdentifier{
			Owner: fork.GetOwner().GetLogin(),
			Repo:  fork.GetName(),
		},
		DefaultBranch: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{
				Owner: fork.GetOwner().GetLogin(),
				Repo:  fork.GetName(),
			},
			Branch: fork.GetDefaultBranch(),
		},
	}, nil
}
