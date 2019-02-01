package usecases

import (
	"context"
	"fmt"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"golang.org/x/sync/errgroup"
)

// SendUpdate provides a use case to send a pull request for updating Gradle in a repository.
type SendUpdate struct {
	dig.In
	GradleService                gateways.GradleService
	RepositoryRepository         gateways.RepositoryRepository
	RepositoryLastScanRepository gateways.RepositoryLastScanRepository
	SendPullRequest              usecases.SendPullRequest
	TimeService                  gateways.TimeService
}

func (usecase *SendUpdate) Do(ctx context.Context, id domain.RepositoryID) error {
	scan := domain.RepositoryLastScan{
		Repository:   id,
		LastScanTime: usecase.TimeService.Now(),
	}
	err := usecase.sendUpdate(ctx, id)
	if err != nil {
		if err, ok := errors.Cause(err).(usecases.SendUpdateError); ok {
			if out := err.PreconditionViolation(); ok {
				scan.PreconditionOut = out
			}
		}
	}
	if err := usecase.RepositoryLastScanRepository.Save(ctx, scan); err != nil {
		return errors.Wrapf(err, "error while saving the scan for the repository %s", id)
	}
	return errors.Wrapf(err, "error while scanning the repository %s", id)
}

func (usecase *SendUpdate) sendUpdate(ctx context.Context, id domain.RepositoryID) error {
	var in domain.GradleUpdatePreconditionIn
	in.BadgeURL = fmt.Sprintf("/%s/%s/status.svg", id.Owner, id.Name) //TODO: externalize

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
		in.Readme = readme
		return nil
	})
	eg.Go(func() error {
		gradleWrapperProperties, err := usecase.RepositoryRepository.GetFileContent(ctx, id, domain.GradleWrapperPropertiesPath)
		if err != nil {
			if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
				if err.NoSuchEntity() {
					return nil
				}
			}
			return errors.Wrapf(err, "error while getting gradle-wrapper.properties")
		}
		in.GradleWrapperProperties = gradleWrapperProperties
		return nil
	})
	eg.Go(func() error {
		latestRelease, err := usecase.GradleService.GetCurrentRelease(ctx)
		if err != nil {
			return errors.Wrapf(err, "error while getting the latest Gradle release")
		}
		in.LatestGradleRelease = latestRelease
		return nil
	})
	if err := eg.Wait(); err != nil {
		return errors.WithStack(err)
	}

	out := domain.CheckGradleUpdatePrecondition(in)
	if out != domain.ReadyToUpdate {
		return errors.WithStack(&sendUpdateError{error: fmt.Errorf("precondition violation (%v)", out), GradleUpdatePreconditionOut: out})
	}

	newProps := domain.ReplaceGradleWrapperVersion(in.GradleWrapperProperties, in.LatestGradleRelease.Version)
	req := usecases.SendPullRequestRequest{
		Base:           id,
		HeadBranchName: fmt.Sprintf("gradle-%s-%s", in.LatestGradleRelease.Version, id.Owner),
		CommitMessage:  fmt.Sprintf("Gradle %s", in.LatestGradleRelease.Version),
		CommitFiles: []domain.File{
			{
				Path:    domain.GradleWrapperPropertiesPath,
				Content: domain.FileContent(newProps),
			},
		},
		Title: fmt.Sprintf("Gradle %s", in.LatestGradleRelease.Version),
		Body:  fmt.Sprintf(`Gradle %s is available.`, in.LatestGradleRelease.Version),
	}
	if err := usecase.SendPullRequest.Do(ctx, req); err != nil {
		return errors.Wrapf(err, "error while sending a pull request %+v", req)
	}
	return nil
}

type sendUpdateError struct {
	error
	GradleUpdatePreconditionOut domain.GradleUpdatePreconditionOut
}

func (err *sendUpdateError) PreconditionViolation() domain.GradleUpdatePreconditionOut {
	return err.GradleUpdatePreconditionOut
}
