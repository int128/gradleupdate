package infrastructure

import (
	"context"
	"os"

	"github.com/google/go-github/v18/github"
	"github.com/gregjones/httpcache"
	"github.com/int128/gradleupdate/infrastructure/memcache"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/log"
)

// GitHubClient creates a client.
func GitHubClient(ctx context.Context) *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Warningf(ctx, "GITHUB_TOKEN is not set")
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	oauth2Client := oauth2.NewClient(ctx, tokenSource)
	cachedTransport := httpcache.Transport{
		Transport: oauth2Client.Transport,
		Cache: memcache.New(ctx),
	}
	return github.NewClient(cachedTransport.Client())
}
