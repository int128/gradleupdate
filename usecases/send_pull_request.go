package usecases

import (
	"context"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
)

type SendPullRequest struct {
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
func (usecase *SendPullRequest) Do(ctx context.Context, req usecases.SendPullRequestRequest) error {
	base, err := usecase.RepositoryRepository.Get(ctx, req.Base)
	if err != nil {
		return errors.Wrapf(err, "could not get the repository %s", req.Base)
	}
	head, err := usecase.RepositoryRepository.Fork(ctx, req.Base)
	if err != nil {
		return errors.Wrapf(err, "could not fork the repository %s", req.Base)
	}
	baseBranch, err := usecase.RepositoryRepository.GetBranch(ctx, base.DefaultBranch)
	if err != nil {
		return errors.Wrapf(err, "could not get the base branch %s", base.DefaultBranch)
	}

	headBranchID := domain.BranchID{
		Repository: head.ID,
		Name:       req.HeadBranchName,
	}
	pushBranchRequest := gateways.PushBranchRequest{
		BaseBranch:    *baseBranch,
		HeadBranch:    headBranchID,
		CommitMessage: req.CommitMessage,
		CommitFiles:   req.CommitFiles,
	}
	headBranch, err := usecase.RepositoryRepository.GetBranch(ctx, headBranchID)
	switch {
	case err == nil:
		if !headBranch.Commit.IsBasedOn(baseBranch.Commit.ID) {
			_, err := usecase.GitService.UpdateForceBranch(ctx, pushBranchRequest)
			if err != nil {
				return errors.Wrapf(err, "could not push the commit to the repository %s", head)
			}
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
		ID:         domain.PullRequestID{Repository: req.Base},
		BaseBranch: base.DefaultBranch,
		HeadBranch: headBranchID,
		Title:      req.Title,
		Body:       req.Body,
	}
	exist, err := usecase.PullRequestRepository.FindByBranch(ctx, pull.BaseBranch, pull.HeadBranch)
	if err != nil {
		return errors.Wrapf(err, "could not find existent pull request")
	}
	if exist != nil {
		return nil // pull request already exists
	}
	if _, err := usecase.PullRequestRepository.Create(ctx, pull); err != nil {
		return errors.Wrapf(err, "could not create a pull request %s", pull)
	}
	return nil
}
