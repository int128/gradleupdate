package usecases

import (
	"context"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/repositories"
	"github.com/pkg/errors"
)

type RepositoryStatus struct {
	Badge      Badge
	Repository domain.Repository
}

type GetRepositoryStatus struct {
	Repository repositories.Repository
}

func (interactor *GetRepositoryStatus) Do(ctx context.Context, owner, repo string) (*RepositoryStatus, error) {
	repository, err := interactor.Repository.Get(ctx, domain.RepositoryIdentifier{Owner: owner, Name: repo})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get the repository %s/%s", owner, repo)
	}
	targetVersion, err := interactor.getVersion(ctx, domain.RepositoryIdentifier{owner, repo})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get version of %s/%s", owner, repo)
	}
	latestVersion, err := interactor.getVersion(ctx, domain.RepositoryIdentifier{Owner: "int128", Name: "latest-gradle-wrapper"})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get the latest version")
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

func (interactor *GetRepositoryStatus) getVersion(ctx context.Context, id domain.RepositoryIdentifier) (domain.GradleVersion, error) {
	file, err := interactor.Repository.GetFile(ctx, id, gradleWrapperPropertiesPath)
	if err != nil {
		return "", errors.Wrapf(err, "File not found: %s", gradleWrapperPropertiesPath)
	}
	v := domain.FindGradleWrapperVersion(string(file.Content))
	if v == "" {
		return "", errors.Errorf("Could not find version from %s", gradleWrapperPropertiesPath)
	}
	return v, nil
}
