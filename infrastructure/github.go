package infrastructure

import (
	"context"
	"os"

	"google.golang.org/appengine/log"

	"github.com/google/go-github/v18/github"
	"golang.org/x/oauth2"
)

// GitHubClient creates a client.
func GitHubClient(ctx context.Context) *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Warningf(ctx, "GITHUB_TOKEN is not set")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
