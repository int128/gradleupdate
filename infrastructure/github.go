package infrastructure

import (
	"net/http"
	"os"

	"github.com/google/go-github/v18/github"
	"go.uber.org/dig"
	"golang.org/x/oauth2"
)

type GitHubClientFactory struct {
	dig.In
	Client *http.Client
}

func (factory *GitHubClientFactory) New() *github.Client {
	//TODO: extract token provider interface
	token := os.Getenv("GITHUB_TOKEN")

	var transport http.RoundTripper
	transport = factory.Client.Transport
	transport = &oauth2.Transport{Base: transport, Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})}
	return github.NewClient(&http.Client{Transport: transport})
}
