package usecases

import (
	"context"
	"fmt"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// SendUpdate provides a use case to send a pull request for updating Gradle in a repository.
type SendUpdate struct {
	dig.In
	SendUpdateQuery         gateways.SendUpdateQuery
	GradleReleaseRepository gateways.GradleReleaseRepository
	PullRequestRepository   gateways.PullRequestRepository
	Logger                  gateways.Logger
}

func (usecase *SendUpdate) Do(ctx context.Context, id git.RepositoryID) error {
	release, err := usecase.GradleReleaseRepository.GetCurrent(ctx)
	if err != nil {
		return errors.Wrapf(err, "error while getting the latest Gradle release")
	}
	out, err := usecase.SendUpdateQuery.Get(ctx, gateways.SendUpdateQueryIn{
		Repository:     id,
		HeadBranchName: gradleupdate.BranchFor(id.Owner, release.Version),
	})
	if err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			if err.NoSuchEntity() {
				return errors.Wrapf(&getRepositoryError{error: err, noSuchRepository: true}, "repository %s not found", id)
			}
		}
		return errors.Wrapf(err, "error while getting the repository %s", id)
	}

	precondition := gradleupdate.Precondition{
		BadgeURL:                gradleupdate.NewBadgeURL(id),
		LatestGradleRelease:     release,
		Readme:                  out.Readme,
		GradleWrapperProperties: out.GradleWrapperProperties,
	}
	preconditionViolation := gradleupdate.CheckPrecondition(precondition)
	if preconditionViolation != gradleupdate.ReadyToUpdate {
		usecase.Logger.Infof(ctx, "skip the repository %v due to precondition violation (%v)", out.BaseRepository, preconditionViolation)
		return errors.WithStack(&sendUpdateError{
			error:                 fmt.Errorf("precondition violation (%v)", preconditionViolation),
			preconditionViolation: preconditionViolation,
		})
	}

	if out.HeadBranch == nil {
		if err := usecase.createBranchAndPullRequest(ctx, out, release); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
	if out.BaseCommitSHA == out.HeadParentCommitSHA {
		usecase.Logger.Infof(ctx, "skip the repository %v because the head branch is up-to-date", out.BaseRepository)
		return nil
	}
	if err := usecase.updateBranch(ctx, out, release); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (usecase *SendUpdate) createBranchAndPullRequest(ctx context.Context, out *gateways.SendUpdateQueryOut, release *gradle.Release) error {
	fork, err := usecase.SendUpdateQuery.ForkRepository(ctx, out.BaseRepository)
	if err != nil {
		return errors.Wrapf(err, "error while forking the repository %s", out.BaseRepository)
	}

	headBranch := gateways.NewBranch{
		Branch: git.BranchID{
			Repository: *fork,
			Name:       gradleupdate.BranchFor(out.BaseRepository.Owner, release.Version),
		},
		ParentCommitSHA: out.BaseCommitSHA,
		ParentTreeSHA:   out.BaseTreeSHA,
		CommitMessage:   formatCommitMessage(release),
		CommitFiles:     formatCommitFiles(out, release),
	}
	if err := usecase.SendUpdateQuery.CreateBranch(ctx, headBranch); err != nil {
		return errors.Wrapf(err, "error while creating a head branch %s", headBranch.Branch)
	}

	pull := git.PullRequest{
		ID:         git.PullRequestID{Repository: out.BaseRepository},
		BaseBranch: out.BaseBranch,
		HeadBranch: headBranch.Branch,
		Title:      formatPullRequestTitle(release),
		Body:       formatPullRequestBody(out.BaseRepository, release),
	}
	if _, err := usecase.PullRequestRepository.Create(ctx, pull); err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			if err.AlreadyExists() {
				usecase.Logger.Infof(ctx, "skip creating a pull request: %s", err)
				return nil
			}
		}
		return errors.Wrapf(err, "error while creating a pull request %s", pull)
	}
	return nil
}

func (usecase *SendUpdate) updateBranch(ctx context.Context, out *gateways.SendUpdateQueryOut, release *gradle.Release) error {
	headBranch := gateways.NewBranch{
		Branch:          *out.HeadBranch,
		ParentCommitSHA: out.BaseCommitSHA,
		ParentTreeSHA:   out.BaseTreeSHA,
		CommitMessage:   formatCommitMessage(release),
		CommitFiles:     formatCommitFiles(out, release),
	}
	if err := usecase.SendUpdateQuery.UpdateBranch(ctx, headBranch, true); err != nil {
		return errors.Wrapf(err, "error while updating the head branch %s", out.HeadBranch)
	}
	return nil
}

func formatCommitMessage(latest *gradle.Release) string {
	return fmt.Sprintf("Gradle %s", latest.Version)
}

func formatCommitFiles(out *gateways.SendUpdateQueryOut, latest *gradle.Release) []git.File {
	newProps := gradle.ReplaceWrapperVersion(out.GradleWrapperProperties, latest.Version)
	return []git.File{
		{
			Path:    gradle.WrapperPropertiesPath,
			Content: git.FileContent(newProps),
		},
	}
}

func formatPullRequestTitle(latest *gradle.Release) string {
	return fmt.Sprintf("Gradle %s", latest.Version)
}

func formatPullRequestBody(id git.RepositoryID, latest *gradle.Release) string {
	return fmt.Sprintf(`Gradle %s is available.

This is sent by @gradleupdate. See %s for more.`,
		latest.Version,
		gradleupdate.NewRepositoryURL(id))
}

type sendUpdateError struct {
	error
	preconditionViolation gradleupdate.PreconditionViolation
}

func (err *sendUpdateError) PreconditionViolation() gradleupdate.PreconditionViolation {
	return err.preconditionViolation
}

var _ usecases.SendUpdateError = &sendUpdateError{}
