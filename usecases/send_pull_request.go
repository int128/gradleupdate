package usecases

import (
	"context"
	"fmt"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/gateways"
	"github.com/pkg/errors"
)

type SendPullRequest struct {
	GradleService         gateways.GradleService
	RepositoryRepository  gateways.RepositoryRepository
	PullRequestRepository gateways.PullRequestRepository
	GitService            gateways.GitService
}

// Do performs:
//
// - Get gradle-wrapper.properties in the repository.
// - Create a file with replaced version.
// - Fork the repository and create a branch with the new file.
// - Create a pull request for the branch.
//
func (usecase *SendPullRequest) Do(ctx context.Context, id domain.RepositoryIdentifier) error {
	latestVersion, err := usecase.GradleService.GetCurrentVersion(ctx)
	if err != nil {
		return errors.Wrapf(err, "could not get the latest Gradle version")
	}
	props, err := usecase.RepositoryRepository.GetFileContent(ctx, id, domain.GradleWrapperPropertiesPath)
	if err != nil {
		return errors.Wrapf(err, "could not find properties file")
	}
	currentVersion := domain.FindGradleWrapperVersion(props.String())
	if currentVersion == "" {
		return errors.Errorf("could not find version in the properties")
	}
	if currentVersion == latestVersion {
		return nil // already up-to-date
	}
	newProps := domain.ReplaceGradleWrapperVersion(props.String(), latestVersion)

	base, err := usecase.RepositoryRepository.Get(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "could not get the repository %s", id)
	}
	head, err := usecase.GitService.ForkBranch(ctx, gateways.ForkBranchRequest{
		Base:           base.DefaultBranch,
		HeadBranchName: fmt.Sprintf("gradle-%s-%s", latestVersion, id.Owner),
		CommitMessage:  fmt.Sprintf("Gradle %s", latestVersion),
		Files: []domain.File{
			{
				Path:    domain.GradleWrapperPropertiesPath,
				Content: []byte(newProps),
			},
		},
	})

	//TODO: if the pull request already exists?
	pull := domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{Repository: id},
		HeadBranch:            head.BranchIdentifier,
		BaseBranch:            base.DefaultBranch,
		Title:                 fmt.Sprintf("Gradle %s", latestVersion),
		Body: fmt.Sprintf(`This will upgrade the Gradle wrapper to the latest version %s.

This pull request is sent by @gradleupdate and based on [int128/latest-gradle-wrapper](https://github.com/int128/latest-gradle-wrapper).
`, latestVersion),
	}
	if _, err := usecase.PullRequestRepository.Create(ctx, pull); err != nil {
		return errors.Wrapf(err, "could not open a pull request %s", pull.String())
	}
	return nil
}
