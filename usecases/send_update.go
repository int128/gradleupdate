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
	"golang.org/x/sync/errgroup"
)

// SendUpdate provides a use case to send a pull request for updating Gradle in a repository.
type SendUpdate struct {
	dig.In
	GradleReleaseRepository gateways.GradleReleaseRepository
	RepositoryRepository    gateways.RepositoryRepository
	SendPullRequest         usecases.SendPullRequest
	Time                    gateways.Time
}

func (usecase *SendUpdate) Do(ctx context.Context, id git.RepositoryID) error {
	precondition := gradleupdate.Precondition{
		BadgeURL: gradleupdate.NewBadgeURL(id),
	}
	var eg errgroup.Group
	eg.Go(func() error {
		readme, err := usecase.RepositoryRepository.GetReadme(ctx, id)
		if err != nil {
			if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
				if err.NoSuchEntity() {
					return nil
				}
			}
			return errors.Wrapf(err, "error while getting README")
		}
		precondition.Readme = readme
		return nil
	})
	eg.Go(func() error {
		gradleWrapperProperties, err := usecase.RepositoryRepository.GetFileContent(ctx, id, gradle.WrapperPropertiesPath)
		if err != nil {
			if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
				if err.NoSuchEntity() {
					return nil
				}
			}
			return errors.Wrapf(err, "error while getting gradle-wrapper.properties")
		}
		precondition.GradleWrapperProperties = gradleWrapperProperties
		return nil
	})
	eg.Go(func() error {
		latestRelease, err := usecase.GradleReleaseRepository.GetCurrent(ctx)
		if err != nil {
			return errors.Wrapf(err, "error while getting the latest Gradle release")
		}
		precondition.LatestGradleRelease = latestRelease
		return nil
	})
	if err := eg.Wait(); err != nil {
		return errors.WithStack(err)
	}

	preconditionViolation := gradleupdate.CheckPrecondition(precondition)
	if preconditionViolation != gradleupdate.ReadyToUpdate {
		return errors.WithStack(&sendUpdateError{error: fmt.Errorf("precondition violation (%v)", preconditionViolation), preconditionViolation: preconditionViolation})
	}

	newProps := gradle.ReplaceWrapperVersion(precondition.GradleWrapperProperties, precondition.LatestGradleRelease.Version)
	req := usecases.SendPullRequestRequest{
		Base:           id,
		HeadBranchName: gradleupdate.BranchFor(id.Owner, precondition.LatestGradleRelease.Version),
		CommitMessage:  fmt.Sprintf("Gradle %s", precondition.LatestGradleRelease.Version),
		CommitFiles: []git.File{
			{
				Path:    gradle.WrapperPropertiesPath,
				Content: git.FileContent(newProps),
			},
		},
		Title: fmt.Sprintf("Gradle %s", precondition.LatestGradleRelease.Version),
		Body: fmt.Sprintf(`Gradle %s is available.

This is sent by @gradleupdate. See %s for more.`,
			precondition.LatestGradleRelease.Version,
			gradleupdate.NewRepositoryURL(id)),
	}
	if err := usecase.SendPullRequest.Do(ctx, req); err != nil {
		return errors.Wrapf(err, "error while sending a pull request %+v", req)
	}
	return nil
}

type sendUpdateError struct {
	error
	preconditionViolation gradleupdate.PreconditionViolation
}

func (err *sendUpdateError) PreconditionViolation() gradleupdate.PreconditionViolation {
	return err.preconditionViolation
}

var _ usecases.SendUpdateError = &sendUpdateError{}
