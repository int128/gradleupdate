package gateways

import (
	"context"
	"testing"

	"github.com/go-test/deep"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/pkg/errors"
)

var sandboxRepository = git.RepositoryID{Owner: "octocat", Name: "Spoon-Knife"}

func TestGetRepositoryQuery_Do_GitHubIntegration(t *testing.T) {
	query := &GetRepositoryQuery{
		Client: gatewaysTestDoubles.NewGitHubClientV4(t),
	}
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		out, err := query.Do(ctx, gateways.GetRepositoryQueryIn{
			Repository:     sandboxRepository,
			HeadBranchName: "example",
		})
		if err != nil {
			t.Fatalf("error while finding the pull request: %+v", err)
		}
		{
			want := git.Repository{
				ID:          sandboxRepository,
				Description: "This repo is for demonstration purposes only.",
				AvatarURL:   "https://avatars3.githubusercontent.com/u/583231?v=4",
				URL:         "https://github.com/octocat/Spoon-Knife",
			}
			if diff := deep.Equal(want, out.Repository); diff != nil {
				t.Error(diff)
			}
		}
		{
			want := git.PullRequestURL("https://github.com/octocat/Spoon-Knife/pull/2180")
			if out.PullRequestURL != want {
				t.Errorf("PullRequestURL wants %s but %s", want, out.PullRequestURL)
			}
		}
		if out.Readme == nil {
			t.Errorf("Readme wants non-nil but nil")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := query.Do(ctx, gateways.GetRepositoryQueryIn{
			Repository:     git.RepositoryID{Owner: "octocat", Name: "not-exist-repository"},
			HeadBranchName: "example",
		})
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		repositoryErr, ok := errors.Cause(err).(gateways.RepositoryError)
		if !ok {
			t.Fatalf("error wants gateways.RepositoryError but %T: %+v", errors.Cause(err), err)
		}
		if !repositoryErr.NoSuchEntity() {
			t.Errorf("NoSuchEntity wants true but false")
		}
	})
}
