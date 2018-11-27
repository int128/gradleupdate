package infrastructure

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/v18/github"
	"github.com/gregjones/httpcache"
	"github.com/int128/gradleupdate/app/infrastructure/memcache"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

// GitHubClient creates a client.
func GitHubClient(ctx context.Context) *github.Client {
	var transport http.RoundTripper
	transport = &urlfetch.Transport{Context: ctx}
	transport = &loggingTransport{ctx, transport}
	transport = &oauth2.Transport{Source: oauth2TokenSource(ctx), Base: transport}
	transport = &httpcache.Transport{Cache: memcache.New(ctx), Transport: transport}
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
