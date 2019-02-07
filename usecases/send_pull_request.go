package usecases

import (
	"context"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type SendPullRequest struct {
	dig.In
	RepositoryRepository  gateways.RepositoryRepository
	PullRequestRepository gateways.PullRequestRepository
	GitService            gateways.GitService
	Logger                gateways.Logger
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
	baseRepository, err := usecase.RepositoryRepository.Get(ctx, req.Base)
	if err != nil {
		return errors.Wrapf(err, "error while getting the base repository %s", req.Base)
	}
	headRepository, err := usecase.RepositoryRepository.Fork(ctx, baseRepository.ID)
	if err != nil {
		return errors.Wrapf(err, "error while forking the repository %s", baseRepository.ID)
	}
	baseBranch, err := usecase.RepositoryRepository.GetBranch(ctx, baseRepository.DefaultBranch)
	if err != nil {
		return errors.Wrapf(err, "error while getting the base branch %s", baseRepository.DefaultBranch)
	}

	headBranchID := git.BranchID{
		Repository: headRepository.ID,
		Name:       req.HeadBranchName,
	}
	pushBranchRequest := gateways.PushBranchRequest{
		BaseBranch:    *baseBranch,
		HeadBranch:    headBranchID,
		CommitMessage: req.CommitMessage,
		CommitFiles:   req.CommitFiles,
	}
	if err := usecase.pushBranch(ctx, pushBranchRequest); err != nil {
		return errors.Wrapf(err, "error while pushing the branch")
	}

	pull := git.PullRequest{
		ID:         git.PullRequestID{Repository: baseRepository.ID},
		BaseBranch: baseRepository.DefaultBranch,
		HeadBranch: headBranchID,
		Title:      req.Title,
		Body:       req.Body,
	}
	if _, err := usecase.PullRequestRepository.Create(ctx, pull); err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			if err.AlreadyExists() {
				usecase.Logger.Infof(ctx, "skip: %s", err)
				return nil
			}
		}
		return errors.Wrapf(err, "error while creating a pull request %s", pull)
	}
	return nil
}

func (usecase *SendPullRequest) pushBranch(ctx context.Context, req gateways.PushBranchRequest) error {
	headBranch, err := usecase.RepositoryRepository.GetBranch(ctx, req.HeadBranch)
	if err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			if err.NoSuchEntity() {
				_, err := usecase.GitService.CreateBranch(ctx, req)
				if err != nil {
					return errors.Wrapf(err, "error while creating a branch %s", req.HeadBranch)
				}
				return nil
			}
		}
		return errors.Wrapf(err, "error while getting the head branch %s", req.HeadBranch)
	}
	if headBranch.Commit.IsBasedOn(req.BaseBranch.Commit.ID) {
		return nil // head branch is up-to-date
	}
	if _, err := usecase.GitService.UpdateForceBranch(ctx, req); err != nil {
		return errors.Wrapf(err, "error while updating the branch %s", req.HeadBranch)
	}
	return nil
}
