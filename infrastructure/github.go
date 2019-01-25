package infrastructure

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/infrastructure/interfaces"
	"go.uber.org/dig"
	"golang.org/x/oauth2"
)

type GitHubClientFactory struct {
	dig.In
	HTTPClientFactory infrastructure.HTTPClientFactory
}

func (factory *GitHubClientFactory) New(ctx context.Context) *github.Client {
	//TODO: extract token provider interface
	token := os.Getenv("GITHUB_TOKEN")

	var transport http.RoundTripper
	transport = factory.HTTPClientFactory.New(ctx).Transport
	transport = &oauth2.Transport{Base: transport, Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})}
	return github.NewClient(&http.Client{Transport: transport})
}
