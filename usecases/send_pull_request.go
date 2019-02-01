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

	headBranchID := git.BranchID{
		Repository: head.ID,
		Name:       req.HeadBranchName,
	}
	pushBranchRequest := gateways.PushBranchRequest{
		BaseBranch:    *baseBranch,
		HeadBranch:    headBranchID,
		CommitMessage: req.CommitMessage,
		CommitFiles:   req.CommitFiles,
	}
	if err := usecase.pushBranch(ctx, pushBranchRequest); err != nil {
		return errors.Wrapf(err, "could not push the branch")
	}

	pull := git.PullRequest{
		ID:         git.PullRequestID{Repository: req.Base},
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

func (usecase *SendPullRequest) pushBranch(ctx context.Context, req gateways.PushBranchRequest) error {
	headBranch, err := usecase.RepositoryRepository.GetBranch(ctx, req.HeadBranch)
	if err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			switch {
			case err.NoSuchEntity():
				_, err := usecase.GitService.CreateBranch(ctx, req)
				if err != nil {
					return errors.Wrapf(err, "could not create a branch %s", req.HeadBranch)
				}
				return nil
			}
		}
		return errors.Wrapf(err, "could not get the head branch %s", req.HeadBranch)
	} else {
		if !headBranch.Commit.IsBasedOn(req.BaseBranch.Commit.ID) {
			_, err := usecase.GitService.UpdateForceBranch(ctx, req)
			if err != nil {
				return errors.Wrapf(err, "could not update the branch %s", req.HeadBranch)
			}
		}
	}
	return nil
}
