package service

import (
	"context"
	"fmt"

	"github.com/int128/gradleupdate/infrastructure"
	"github.com/int128/gradleupdate/service/pr"
	"github.com/pkg/errors"
)

func CreateOrUpdatePullRequestForGradleWrapper(ctx context.Context, owner, repo, version string) error {
	c := infrastructure.GitHubClient(ctx)
	baseRepository, _, err := c.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return errors.Wrapf(err, "Could not get the repository %s/%s", owner, repo)
	}
	headRepository, err := pr.Fork(ctx, c, pr.Repository{Owner: owner, Repo: repo})
	if err != nil {
		return errors.Wrapf(err, "Could not fork the repository %s/%s", owner, repo)
	}

	base := pr.Branch{
		Repository: pr.Repository{Owner: owner, Repo: repo},
		Branch:     baseRepository.GetDefaultBranch(),
	}
	head := pr.Branch{
		Repository: pr.Repository{
			Owner: headRepository.GetOwner().GetLogin(),
			Repo:  headRepository.GetName(),
		},
		Branch: fmt.Sprintf("gradle-%s-%s", version, owner),
	}
	commit := pr.Commit{
		Message: fmt.Sprintf("Gradle %s", version),
		Files: []pr.File{
			// TODO: Gradle wrapper files
			pr.File{
				Path:           "example",
				Mode:           "100644",
				EncodedContent: "MjAxOOW5tCAxMeaciDEz5pelIOeBq+abnOaXpSAxN+aZgjEy5YiGMDXnp5IgSlNUCg==",
			},
		},
	}
	if err := pr.CreateOrUpdateBranch(ctx, c, base, head, commit); err != nil {
		return errors.Wrapf(err, "Could not create or update the branch")
	}

	pull := pr.PullRequest{
		Head:  head,
		Base:  base,
		Title: fmt.Sprintf("Gradle %s", version),
	}
	if _, err := pr.CreateOrUpdatePullRequest(ctx, c, pull); err != nil {
		return errors.Wrapf(err, "Could not create or update the pull request")
	}
	return nil
}
