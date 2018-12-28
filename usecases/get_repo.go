package usecases

import (
	"context"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
)

type GetRepositoryResponse struct {
	Repository     domain.Repository
	CurrentVersion domain.GradleVersion
	LatestVersion  domain.GradleVersion
	UpToDate       bool
}

type GetRepository struct {
	GradleService        gateways.GradleService
	RepositoryRepository gateways.RepositoryRepository
}

func (usecase *GetRepository) Do(ctx context.Context, id domain.RepositoryID) (*GetRepositoryResponse, error) {
	repository, err := usecase.RepositoryRepository.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get the repository %s", id)
	}

	props, err := usecase.RepositoryRepository.GetFileContent(ctx, id, domain.GradleWrapperPropertiesPath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get the properties file in %s", id)
	}
	currentVersion := domain.FindGradleWrapperVersion(props.String())
	if currentVersion == "" {
		return nil, errors.Errorf("could not find version from properties file in %s", id)
	}
	latestVersion, err := usecase.GradleService.GetCurrentVersion(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get the latest Gradle version")
	}
	return &GetRepositoryResponse{
		Repository:     *repository,
		CurrentVersion: currentVersion,
		LatestVersion:  latestVersion,
		UpToDate:       currentVersion.GreaterOrEqualThan(latestVersion),
	}, nil
}
