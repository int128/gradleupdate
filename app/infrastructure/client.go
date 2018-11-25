package infrastructure

import (
	"context"
	"google.golang.org/appengine/urlfetch"
	"net/http"
	"os"

	"github.com/google/go-github/v18/github"
	"github.com/gregjones/httpcache"
	"github.com/int128/gradleupdate/app/infrastructure/memcache"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/log"
)

// GitHubClient creates a client.
func GitHubClient(ctx context.Context) *github.Client {
	appengineTransport := &urlfetch.Transport{Context: ctx}

	loggingTransport := &loggingTransport{ctx, appengineTransport}

	oauth2Transport := oauth2Transport(ctx, loggingTransport)

	cachedTransport := &httpcache.Transport{Cache: memcache.New(ctx), Transport: oauth2Transport}

	return github.NewClient(&http.Client{Transport: cachedTransport})
}

func oauth2Transport(ctx context.Context, base http.RoundTripper) http.RoundTripper {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Warningf(ctx, "GITHUB_TOKEN is not set")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return &oauth2.Transport{
		Base:   base,
		Source: ts,
	}
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
