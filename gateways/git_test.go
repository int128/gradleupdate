package gateways

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"golang.org/x/oauth2"
)

var sandboxRepository = domain.RepositoryID{Owner: "octocat", Name: "Spoon-Knife"}

func TestGitService_CreateBranch(t *testing.T) {
	client := newGitHubClient(t)
	service := GitService{GitHubClientFactory: &factory{client}}
	ctx := context.Background()
	fork, _, err := client.Repositories.CreateFork(ctx, sandboxRepository.Owner, sandboxRepository.Name, nil)
	if err != nil {
		if _, ok := err.(*github.AcceptedError); !ok {
			t.Fatalf("could not fork the repository %s: %s", sandboxRepository, err)
		}
	}
	base, _, err := client.Repositories.GetBranch(ctx, sandboxRepository.Owner, sandboxRepository.Name, "master")
	if err != nil {
		t.Fatalf("could not get the master branch: %s", err)
	}
	req := gateways.PushBranchRequest{
		BaseBranch: domain.Branch{
			ID: domain.BranchID{
				Repository: sandboxRepository,
				Name:       base.GetName(),
			},
			Commit: domain.Commit{
				ID: domain.CommitID{
					Repository: sandboxRepository,
					SHA:        domain.CommitSHA(base.Commit.GetSHA()),
				},
				Parents: []domain.CommitID{ /* omit */ },
				Tree: domain.TreeID{
					Repository: sandboxRepository,
					SHA:        domain.TreeSHA(base.Commit.Commit.Tree.GetSHA()),
				},
			},
		},
		HeadBranch: domain.BranchID{
			Repository: domain.RepositoryID{Owner: fork.GetOwner().GetLogin(), Name: fork.GetName()},
			Name:       "branch1",
		},
		CommitMessage: "Test commit",
		CommitFiles: []domain.File{
			{
				Path:    "foo/bar",
				Content: domain.FileContent("baz"),
			},
		},
	}

	t.Run("DeleteHeadBranchBeforeTest", func(t *testing.T) {
		if resp, err := client.Git.DeleteRef(ctx, fork.GetOwner().GetLogin(), fork.GetName(), req.HeadBranch.Ref()); err != nil {
			if resp.StatusCode != 404 {
				t.Fatalf("could not remove the branch %s", req.HeadBranch)
			}
		}
	})
	t.Run("ShouldCreateBranch", func(t *testing.T) {
		branch, err := service.CreateBranch(ctx, req)
		if err != nil {
			t.Fatalf("error on CreateBranch: %s", err)
		}
		if branch.ID != req.HeadBranch {
			t.Errorf("branch.Name wants %s but %s", req.HeadBranch, branch.ID)
		}
		if branch.Commit.ID.SHA == "" {
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
	resp, err := t.transport.RoundTrip(req)
	if resp != nil {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Printf("could not dump response: %s", err)
		}
		log.Printf("%s %s\n%s", req.Method, req.URL, string(dump))
	}
	return resp, err
}
