package infrastructure

import (
	"net/http"

	"github.com/google/go-github/v24/github"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
	"go.uber.org/dig"
	"golang.org/x/oauth2"
)

type GitHubClientFactory struct {
	dig.In
	Client      *http.Client
	Credentials gateways.Credentials
}

func (factory *GitHubClientFactory) NewV3() *github.Client {
	var transport http.RoundTripper
	transport = factory.Client.Transport
	transport = &oauth2Transport{transport, factory.Credentials}
	return github.NewClient(&http.Client{Transport: transport})
}

func (factory *GitHubClientFactory) NewV4() *githubv4.Client {
	var transport http.RoundTripper
	transport = factory.Client.Transport
	transport = &oauth2Transport{transport, factory.Credentials}
	return githubv4.NewClient(&http.Client{Transport: transport})
}

type oauth2Transport struct {
	Transport   http.RoundTripper
	Credentials gateways.Credentials
}

func (t *oauth2Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	credentials, err := t.Credentials.Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get credentials for GitHub API")
	}
	transport := &oauth2.Transport{Base: t.Transport, Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: credentials.GitHubToken})}
	return transport.RoundTrip(req)
}
