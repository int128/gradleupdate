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

type GetRepository struct {
	dig.In
	GradleService        gateways.GradleService
	RepositoryRepository gateways.RepositoryRepository
}

func (usecase *GetRepository) Do(ctx context.Context, id domain.RepositoryID) (*usecases.GetRepositoryResponse, error) {
	repository, err := usecase.RepositoryRepository.Get(ctx, id)
	if err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			if err.NoSuchEntity() {
				return nil, errors.Wrapf(&getRepositoryError{error: err, noSuchRepository: true}, "repository %s not found", id)
			}
		}
		return nil, errors.Wrapf(err, "error while getting the repository %s", id)
	}

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
		return nil, errors.WithStack(err)
	}

	out := domain.CheckGradleUpdatePrecondition(in)
	return &usecases.GetRepositoryResponse{
		Repository:                  *repository,
		GradleUpdatePreconditionOut: out,
	}, nil
}

type getRepositoryError struct {
	error
	noSuchRepository bool
}

func (err *getRepositoryError) NoSuchRepository() bool { return err.noSuchRepository }
