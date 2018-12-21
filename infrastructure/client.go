package infrastructure

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/infrastructure/httpcache"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func NewClient(ctx context.Context) *http.Client {
	var transport http.RoundTripper
	transport = &urlfetch.Transport{Context: ctx}
	transport = &loggingTransport{ctx, transport}
	transport = &httpcache.Transport{Cache: &httpcache.AppEngineMemcache{Context: ctx}, Transport: transport}
	return &http.Client{Transport: transport}
}

// GitHubClient creates a client.
func GitHubClient(ctx context.Context) *github.Client {
	var transport http.RoundTripper
	transport = &urlfetch.Transport{Context: ctx}
	transport = &loggingTransport{ctx, transport}
	transport = &httpcache.Transport{Cache: &httpcache.AppEngineMemcache{Context: ctx}, Transport: transport}
	transport = &oauth2.Transport{Source: oauth2TokenSource(ctx), Base: transport}
	return github.NewClient(&http.Client{Transport: transport})
}

func oauth2TokenSource(ctx context.Context) oauth2.TokenSource {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Warningf(ctx, "GITHUB_TOKEN is not set")
	}
	return oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
}

type loggingTransport struct {
	ctx  context.Context
	base http.RoundTripper
}

func (t loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.base.RoundTrip(req)
	if res != nil {
		log.Debugf(t.ctx, "[GitHubClient] %d %s %s", res.StatusCode, req.Method, req.URL)
	}
	return res, err
}
