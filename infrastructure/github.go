package infrastructure

import (
	"context"
	"net/http"
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
	oauth2Transport := oauth2Client.Transport

	cachedTransport := httpcache.Transport{
		Transport: oauth2Transport,
		Cache:     memcache.New(ctx),
	}

	return github.NewClient(&http.Client{
		Transport: loggingTransport{ctx, &cachedTransport},
	})
}

type loggingTransport struct {
	ctx       context.Context
	transport http.RoundTripper
}

func (t loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.transport.RoundTrip(req)
	if res != nil {
		log.Debugf(t.ctx, "[GitHubClient] %d %s %s", res.StatusCode, req.Method, req.URL)
	}
	return res, err
}
