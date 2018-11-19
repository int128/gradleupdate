package service

import (
	"context"
	"fmt"

	"github.com/google/go-github/v18/github"
)

// GetRepository returns the repository.
func GetRepository(ctx context.Context, owner, repo string) (*github.Repository, error) {
	c := github.Client{} //TODO
	r, _, err := c.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("Could not get the repository: %s", err)
	}
	return r, nil
}
