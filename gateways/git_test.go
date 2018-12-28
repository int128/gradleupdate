package gateways

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"golang.org/x/oauth2"
)

var sandboxRepository = domain.RepositoryIdentifier{Owner: "octocat", Name: "Spoon-Knife"}

func TestGitService_ForkBranch(t *testing.T) {
	client := newGitHubClient(t)
	service := GitService{GitHubClientFactory: &factory{client}}
	ctx := context.Background()
	req := gateways.ForkBranchRequest{
		Base: domain.BranchIdentifier{
			Repository: sandboxRepository,
			Name:       "master",
		},
		HeadBranchName: "branch1",
		CommitMessage:  "Test commit",
		Files: []domain.File{
			{
				Path:    "foo/bar",
				Content: domain.FileContent("baz"),
			},
		},
	}

	t.Run("DeleteHeadBranchBeforeTest", func(t *testing.T) {
		fork, _, err := client.Repositories.CreateFork(ctx, sandboxRepository.Owner, sandboxRepository.Name, nil)
		if err != nil {
			if _, ok := err.(*github.AcceptedError); !ok {
				t.Fatalf("could not fork the repository %s", sandboxRepository)
			}
		}
		if resp, err := client.Git.DeleteRef(ctx, fork.GetOwner().GetLogin(), fork.GetName(), "refs/heads/"+req.HeadBranchName); err != nil {
			if resp.StatusCode != 404 {
				t.Fatalf("could not remove the branch %s", req.HeadBranchName)
			}
		}
	})
	t.Run("ShouldCreateBranch", func(t *testing.T) {
		branch, err := service.ForkBranch(ctx, req)
		if err != nil {
			t.Fatalf("error from ForkBranch: %s", err)
		}
		if branch.Name != req.HeadBranchName {
			t.Errorf("branch.Name wants %s but %s", req.HeadBranchName, branch.Name)
		}
		if branch.CommitSHA == "" {
			t.Errorf("branch.CommitSHA wants non-empty but empty")
		}
	})
	t.Run("ShouldDoNothing", func(t *testing.T) {
		branch, err := service.ForkBranch(ctx, req)
		if err != nil {
			t.Fatalf("error from ForkBranch: %s", err)
		}
		if branch.Name != req.HeadBranchName {
			t.Errorf("branch.Name wants %s but %s", req.HeadBranchName, branch.Name)
		}
		if branch.CommitSHA == "" {
			t.Errorf("branch.CommitSHA wants non-empty but empty")
		}
	})
}

func newGitHubClient(t *testing.T) *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skipf("GITHUB_TOKEN is not set and skip the test")
	}
	var transport http.RoundTripper
	transport = http.DefaultTransport
	transport = &loggingTransport{transport}
	transport = &oauth2.Transport{Base: transport, Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})}
	return github.NewClient(&http.Client{Transport: transport})
}

type factory struct {
	client *github.Client
}

func (c *factory) New(ctx context.Context) *github.Client {
	return c.client
}

type loggingTransport struct {
	transport http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.transport.RoundTrip(req)
	if res != nil {
		log.Printf("%d %s %s", res.StatusCode, req.Method, req.URL)
	}
	return res, err
}
