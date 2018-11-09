package service

import (
	"context"
	"fmt"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/infrastructure"
)

func GetRepository(ctx context.Context, owner, repo string) (*github.Repository, error) {
	c := infrastructure.GitHubClient(ctx)
	r, _, err := c.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("Could not get the repository: %s", err)
	}
	return r, nil
}
