package infrastructure

import (
	"net/http"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"golang.org/x/oauth2"
)

type GitHubClientFactory struct {
	dig.In
	Client           *http.Client
	ConfigRepository gateways.ConfigRepository
}

func (factory *GitHubClientFactory) New() *github.Client {
	var transport http.RoundTripper
	transport = factory.Client.Transport
	transport = &oauth2Transport{transport, factory.ConfigRepository}
	return github.NewClient(&http.Client{Transport: transport})
}

type oauth2Transport struct {
	Transport        http.RoundTripper
	ConfigRepository gateways.ConfigRepository
}

func (t *oauth2Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	config, err := t.ConfigRepository.Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting Config")
	}
	transport := &oauth2.Transport{Base: t.Transport, Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.GitHubToken})}
	return transport.RoundTrip(req)
}
