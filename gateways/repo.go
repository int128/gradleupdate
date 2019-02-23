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
