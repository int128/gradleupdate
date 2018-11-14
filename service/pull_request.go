package service

import (
	"context"
	"fmt"

	"github.com/int128/gradleupdate/infrastructure"
	"github.com/int128/gradleupdate/service/pr"
	"github.com/pkg/errors"
)

// CreateOrUpdatePullRequestForGradleWrapper opens a pull request for updating the wrapper.
func CreateOrUpdatePullRequestForGradleWrapper(ctx context.Context, owner, repo, version string) error {
	c := infrastructure.GitHubClient(ctx)

	files, err := FindGradleWrapperFiles(ctx, c, "int128", "latest-gradle-wrapper")
	if err != nil {
		return errors.Wrapf(err, "Could not find files of the latest Gradle wrapper")
	}

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
		Files:   files,
	}
	if err := pr.CreateOrUpdateBranch(ctx, c, base, head, commit); err != nil {
		return errors.Wrapf(err, "Could not create or update the branch")
	}

	pull := pr.PullRequest{
		Head:  head,
		Base:  base,
		Title: fmt.Sprintf("Gradle %s", version),
		Body:  fmt.Sprintf(`This will upgrade the Gradle wrapper to the latest version %s.

This pull request is sent by @gradleupdate and based on [int128/latest-gradle-wrapper](https://github.com/int128/latest-gradle-wrapper).
`, version),
	}
	if _, err := pr.CreateOrUpdatePullRequest(ctx, c, pull); err != nil {
		return errors.Wrapf(err, "Could not create or update the pull request")
	}
	return nil
}
