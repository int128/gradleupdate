package gateways

import (
	"context"

	"github.com/google/go-github/v18/github"
)

type GitHubClientFactory interface {
	New(ctx context.Context) *github.Client
}
