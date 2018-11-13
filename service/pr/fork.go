package pr

import (
	"context"

	"github.com/google/go-github/v18/github"
	"github.com/pkg/errors"
)

// Fork forks the repository.
func Fork(ctx context.Context, c *github.Client, r Repository) (*github.Repository, error) {
	fork, _, err := c.Repositories.CreateFork(ctx, r.Owner, r.Repo, &github.RepositoryCreateForkOptions{})
	if _, ok := err.(*github.AcceptedError); ok {
		return fork, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "Could not fork the repository %s", r)
	}
	return fork, nil
}
