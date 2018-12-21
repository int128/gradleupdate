package usecases

import (
	"context"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/gateways"
	"github.com/pkg/errors"
)

type RepositoryStatus struct {
	Badge      Badge
	Repository domain.Repository
}

type GetRepository struct {
	GradleService        gateways.GradleService
	RepositoryRepository gateways.RepositoryRepository
}

func (usecase *GetRepository) Do(ctx context.Context, owner, repo string) (*RepositoryStatus, error) {
	repository, err := usecase.RepositoryRepository.Get(ctx, domain.RepositoryIdentifier{Owner: owner, Name: repo})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get the repository %s/%s", owner, repo)
	}
	targetVersion, err := usecase.getVersion(ctx, domain.RepositoryIdentifier{owner, repo})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get version of %s/%s", owner, repo)
	}
	latestVersion, err := usecase.GradleService.GetCurrentVersion(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get the latest Gradle version")
	}
	return &RepositoryStatus{
		Badge: Badge{
			TargetVersion: targetVersion,
			LatestVersion: latestVersion,
			UpToDate:      domain.IsUpToDate(targetVersion, latestVersion),
		},
		Repository: repository,
	}, nil
}

func (usecase *GetRepository) getVersion(ctx context.Context, id domain.RepositoryIdentifier) (domain.GradleVersion, error) {
	file, err := usecase.RepositoryRepository.GetFile(ctx, id, gradleWrapperPropertiesPath)
	if err != nil {
		return "", errors.Wrapf(err, "File not found: %s", gradleWrapperPropertiesPath)
	}
	v := domain.FindGradleWrapperVersion(string(file.Content))
	if v == "" {
		return "", errors.Errorf("Could not find version from %s", gradleWrapperPropertiesPath)
	}
	return v, nil
}
