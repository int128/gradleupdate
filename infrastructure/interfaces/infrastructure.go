package infrastructure

import (
	"context"
	"net/http"

	"github.com/google/go-github/v18/github"
)

type HTTPClientFactory interface {
	New(ctx context.Context) *http.Client
}

type GitHubClientFactory interface {
	New(ctx context.Context) *github.Client
}
