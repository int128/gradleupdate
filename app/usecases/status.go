package usecases

import (
	"context"
	"github.com/int128/gradleupdate/app/domain"
	"github.com/int128/gradleupdate/app/domain/repositories"
	"github.com/pkg/errors"
)

// Status represents whether the wrapper is up-to-date or out-of-date.
type Status struct {
	TargetVersion domain.GradleVersion
	LatestVersion domain.GradleVersion
	UpToDate      bool
}

// GetStatus provides a usecase to get status of Gradle wrapper in a repository.
type GetStatus struct {
	Repository repositories.Repository
}

// Do performs the usecase.
func (interactor *GetStatus) Do(ctx context.Context, owner, repo string) (*Status, error) {
	targetVersion, err := interactor.getVersion(ctx, domain.RepositoryIdentifier{owner, repo})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get version of %s/%s", owner, repo)
	}
	latestVersion, err := interactor.getVersion(ctx, domain.RepositoryIdentifier{Owner: "int128", Name: "latest-gradle-wrapper"})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get the latest version")
	}
	return &Status{
		TargetVersion: targetVersion,
		LatestVersion: latestVersion,
		UpToDate:      domain.IsUpToDate(targetVersion, latestVersion),
	}, nil
}

func (interactor *GetStatus) getVersion(ctx context.Context, id domain.RepositoryIdentifier) (domain.GradleVersion, error) {
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
