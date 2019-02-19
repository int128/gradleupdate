package gatewaysTestDoubles

import (
	"net/http"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/google/go-github/v24/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// NewGitHubClientV3 returns a client for GitHub integration tests.
// This requires env GITHUB_TOKEN and skips the test if it is not set.
func NewGitHubClientV3(t *testing.T) *github.Client {
	t.Helper()
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skipf("GITHUB_TOKEN is not set and skip the test")
	}
	var transport http.RoundTripper
	transport = http.DefaultTransport
	transport = &loggingTransport{t, transport}
	transport = &oauth2.Transport{Base: transport, Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})}
	return github.NewClient(&http.Client{Transport: transport})
}

// NewGitHubClientV4 returns a client for GitHub integration tests.
// This requires env GITHUB_TOKEN and skips the test if it is not set.
func NewGitHubClientV4(t *testing.T) *githubv4.Client {
	t.Helper()
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skipf("GITHUB_TOKEN is not set and skip the test")
	}
	var transport http.RoundTripper
	transport = http.DefaultTransport
	transport = &loggingTransport{t, transport}
	transport = &oauth2.Transport{Base: transport, Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})}
	return githubv4.NewClient(&http.Client{Transport: transport})
}

type loggingTransport struct {
	t         *testing.T
	transport http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.transport.RoundTrip(req)
	if resp != nil {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			t.t.Errorf("could not dump response: %s", err)
		}
		t.t.Logf("%s %s\n%s", req.Method, req.URL, string(dump))
	}
	return resp, err
}
