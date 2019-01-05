package infrastructure

import (
	"context"
	"net/http"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/infrastructure/httpcache"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/urlfetch"
)

type GitHubClientFactory struct {
	Token                   string
	ResponseCacheRepository gateways.ResponseCacheRepository
}

func (c *GitHubClientFactory) New(ctx context.Context) *github.Client {
	var transport http.RoundTripper
	transport = &urlfetch.Transport{Context: ctx}
	transport = &loggingTransport{Transport: transport, Context: ctx, Name: "GitHubClient"}
	transport = &httpcache.Transport{Transport: transport, Context: ctx, ResponseCacheRepository: c.ResponseCacheRepository}
	transport = &oauth2.Transport{Base: transport, Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.Token})}
	return github.NewClient(&http.Client{Transport: transport})
}
