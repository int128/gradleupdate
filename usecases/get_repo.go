package usecases

import (
	"context"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"golang.org/x/sync/errgroup"
)

type GetRepository struct {
	dig.In
	GradleReleaseRepository gateways.GradleReleaseRepository
	RepositoryRepository    gateways.RepositoryRepository
}

func (usecase *GetRepository) Do(ctx context.Context, id git.RepositoryID) (*usecases.GetRepositoryResponse, error) {
	repository, err := usecase.RepositoryRepository.Get(ctx, id)
	if err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			if err.NoSuchEntity() {
				return nil, errors.Wrapf(&getRepositoryError{error: err, noSuchRepository: true}, "repository %s not found", id)
			}
		}
		return nil, errors.Wrapf(err, "error while getting the repository %s", id)
	}

	precondition := gradleupdate.Precondition{
		BadgeURL: gradleupdate.NewBadgeURL(repository.ID),
	}
	var eg errgroup.Group
	eg.Go(func() error {
		readme, err := usecase.RepositoryRepository.GetReadme(ctx, repository.ID)
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
		gradleWrapperProperties, err := usecase.RepositoryRepository.GetFileContent(ctx, repository.ID, gradle.WrapperPropertiesPath)
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
		return nil, errors.WithStack(err)
	}

	preconditionViolation := gradleupdate.CheckPrecondition(precondition)
	return &usecases.GetRepositoryResponse{
		Repository:                  *repository,
		LatestGradleRelease:         *precondition.LatestGradleRelease,
		UpdatePreconditionViolation: preconditionViolation,
	}, nil
}

type getRepositoryError struct {
	error
	noSuchRepository bool
}

func (err *getRepositoryError) NoSuchRepository() bool { return err.noSuchRepository }
