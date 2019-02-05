package usecases

import (
	"context"
	"fmt"

	"github.com/int128/gradleupdate/domain"
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
	GradleService                  gateways.GradleService
	RepositoryRepository           gateways.RepositoryRepository
	RepositoryLastUpdateRepository gateways.RepositoryLastUpdateRepository
	SendPullRequest                usecases.SendPullRequest
	TimeService                    gateways.TimeService
}

func (usecase *SendUpdate) Do(ctx context.Context, id git.RepositoryID) error {
	lastUpdate := domain.RepositoryLastUpdate{
		Repository:     id,
		LastUpdateTime: usecase.TimeService.Now(),
	}
	err := usecase.sendUpdate(ctx, id)
	if err != nil {
		if err, ok := errors.Cause(err).(usecases.SendUpdateError); ok {
			lastUpdate.PreconditionViolation = err.PreconditionViolation()
		}
	}
	if err := usecase.RepositoryLastUpdateRepository.Save(ctx, lastUpdate); err != nil {
		return errors.Wrapf(err, "error while saving the scan for the repository %s", id)
	}
	return errors.Wrapf(err, "error while scanning the repository %s", id)
}

func (usecase *SendUpdate) sendUpdate(ctx context.Context, id git.RepositoryID) error {
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
		latestRelease, err := usecase.GradleService.GetCurrentRelease(ctx)
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
		HeadBranchName: fmt.Sprintf("gradle-%s-%s", precondition.LatestGradleRelease.Version, id.Owner),
		CommitMessage:  fmt.Sprintf("Gradle %s", precondition.LatestGradleRelease.Version),
		CommitFiles: []git.File{
			{
				Path:    gradle.WrapperPropertiesPath,
				Content: git.FileContent(newProps),
			},
		},
		Title: fmt.Sprintf("Gradle %s", precondition.LatestGradleRelease.Version),
		Body:  fmt.Sprintf(`Gradle %s is available.`, precondition.LatestGradleRelease.Version),
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
