package usecases

import (
	"context"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
)

type GetRepository struct {
	GradleService        gateways.GradleService
	RepositoryRepository gateways.RepositoryRepository
}

func (usecase *GetRepository) Do(ctx context.Context, id domain.RepositoryID) (*usecases.GetRepositoryResponse, error) {
	repository, err := usecase.RepositoryRepository.Get(ctx, id)
	if err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			switch {
			case err.NoSuchEntity():
				return nil, errors.Wrapf(&getRepositoryError{error: err, noSuchRepository: true}, "repository %s not found", id)
			}
		}
		return nil, errors.Wrapf(err, "could not get the repository %s", id)
	}

	gradleWrapperProperties, err := usecase.RepositoryRepository.GetFileContent(ctx, id, domain.GradleWrapperPropertiesPath)
	if err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			switch {
			case err.NoSuchEntity():
				return nil, errors.Wrapf(&getRepositoryError{error: err, noGradleVersion: true}, "gradleWrapperProperties not found")
			}
		}
		return nil, errors.Wrapf(err, "could not get the properties file in %s", id)
	}

	currentVersion := domain.FindGradleWrapperVersion(gradleWrapperProperties)
	if currentVersion == "" {
		return nil, errors.Wrapf(&getRepositoryError{noGradleVersion: true}, "could not find version from properties file in %s", id)
	}

	latestVersion, err := usecase.GradleService.GetCurrentVersion(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get the latest Gradle version")
	}

	return &usecases.GetRepositoryResponse{
		Repository:     *repository,
		CurrentVersion: currentVersion,
		LatestVersion:  latestVersion,
		UpToDate:       currentVersion.GreaterOrEqualThan(latestVersion),
	}, nil
}

type getRepositoryError struct {
	error
	noSuchRepository bool
	noGradleVersion  bool
}

func (err *getRepositoryError) NoSuchRepository() bool { return err.noSuchRepository }
func (err *getRepositoryError) NoGradleVersion() bool  { return err.noGradleVersion }
