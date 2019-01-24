package infrastructure

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/infrastructure/httpcache"
	"go.uber.org/dig"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/urlfetch"
)

type GitHubClientFactory struct {
	dig.In
	ResponseCacheRepository gateways.ResponseCacheRepository
	Logger                  gateways.Logger
}

func (c *GitHubClientFactory) New(ctx context.Context) *github.Client {
	//TODO: extract token provider interface
	token := os.Getenv("GITHUB_TOKEN")

	var transport http.RoundTripper
	transport = &urlfetch.Transport{Context: ctx}
	transport = &loggingTransport{Transport: transport, Context: ctx, Name: "GitHubClient", Logger: c.Logger}
	transport = &httpcache.Transport{Transport: transport, Context: ctx, ResponseCacheRepository: c.ResponseCacheRepository, Logger: c.Logger}
	transport = &oauth2.Transport{Base: transport, Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})}
	return github.NewClient(&http.Client{Transport: transport})
}
