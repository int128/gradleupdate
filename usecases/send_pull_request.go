package usecases

import (
	"context"
	"fmt"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
)

type SendPullRequest struct {
	GradleService         gateways.GradleService
	RepositoryRepository  gateways.RepositoryRepository
	PullRequestRepository gateways.PullRequestRepository
	GitService            gateways.GitService
}

// Do sends a pull request for updating the Gradle wrapper in the repository.
//
// If the Gradle wrapper is up-to-date, do nothing.
//
// If the head branch exists, check if its parent is the base branch.
// If not, update the head branch onto the base branch.
//
// If the pull request exists, do not create any more.
//
func (usecase *SendPullRequest) Do(ctx context.Context, id domain.RepositoryID) error {
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
	head, err := usecase.RepositoryRepository.Fork(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "could not fork the repository %s", id)
	}
	baseBranch, err := usecase.RepositoryRepository.GetBranch(ctx, base.DefaultBranch)
	if err != nil {
		return errors.Wrapf(err, "could not get the base branch %s", base.DefaultBranch)
	}

	headBranchID := head.ID.Branch(fmt.Sprintf("gradle-%s-%s", latestVersion, id.Owner))
	pushBranchRequest := gateways.PushBranchRequest{
		BaseBranch:    *baseBranch,
		HeadBranch:    headBranchID,
		CommitMessage: fmt.Sprintf("Gradle %s", latestVersion),
		CommitFiles: []domain.File{
			{
				Path:    domain.GradleWrapperPropertiesPath,
				Content: domain.FileContent(newProps),
			},
		},
	}
	headBranch, err := usecase.RepositoryRepository.GetBranch(ctx, headBranchID)
	switch {
	case err == nil:
		if headBranch.Commit.IsBasedOn(baseBranch.Commit.ID) {
			return nil
		}
		_, err := usecase.GitService.UpdateForceBranch(ctx, pushBranchRequest)
		if err != nil {
			return errors.Wrapf(err, "could not push the commit to the repository %s", head)
		}
	case usecase.RepositoryRepository.IsNotFoundError(err):
		_, err := usecase.GitService.CreateBranch(ctx, pushBranchRequest)
		if err != nil {
			return errors.Wrapf(err, "could not push the commit to the repository %s", head)
		}
	default:
		return errors.Wrapf(err, "could not get the head branch %s", head)
	}

	pull := domain.PullRequest{
		ID:         domain.PullRequestID{Repository: id},
		HeadBranch: headBranchID,
		BaseBranch: base.DefaultBranch,
		Title:      fmt.Sprintf("Gradle %s", latestVersion),
		Body:       fmt.Sprintf(`Gradle %s is available.`, latestVersion),
	}
	if _, err := usecase.PullRequestRepository.Create(ctx, pull); err != nil {
		return errors.Wrapf(err, "could not create a pull request %s", pull)
	}
	return nil
}
