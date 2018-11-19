package usecases

import (
	"context"
	"fmt"

	"github.com/int128/gradleupdate/app/domain"
	"github.com/int128/gradleupdate/app/domain/repositories"
	"github.com/pkg/errors"
)

const gradleWrapperPropsPath = "gradle/wrapper/gradle-wrapper.properties"

// GradleWrapperStatus represents whether the wrapper is up-to-date or out-of-date.
type GradleWrapperStatus struct {
	TargetVersion domain.GradleVersion
	LatestVersion domain.GradleVersion
	UpToDate      bool
}

// GetGradleWrapperStatus provides a usecase to get status of Gradle wrapper in a repository.
type GetGradleWrapperStatus struct {
	Repository repositories.Repository
}

// Do performs the usecase.
func (interactor *GetGradleWrapperStatus) Do(ctx context.Context, id domain.RepositoryIdentifier) (*GradleWrapperStatus, error) {
	targetVersion, err := interactor.getVersion(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get version of %s", id)
	}
	latestVersion, err := interactor.getVersion(ctx, domain.RepositoryIdentifier{
		Owner: "int128", Repo: "latest-gradle-wrapper",
	})
	if err != nil {
		return nil, fmt.Errorf("Could not get the latest version: %s", err)
	}
	return &GradleWrapperStatus{
		TargetVersion: targetVersion,
		LatestVersion: latestVersion,
		UpToDate:      domain.IsUpToDate(targetVersion, latestVersion),
	}, nil
}

func (interactor *GetGradleWrapperStatus) getVersion(ctx context.Context, id domain.RepositoryIdentifier) (domain.GradleVersion, error) {
	file, err := interactor.Repository.GetFile(ctx, id, gradleWrapperPropsPath)
	if err != nil {
		return "", errors.Wrapf(err, "File not found: %s", gradleWrapperPropsPath)
	}
	v := domain.FindGradleWrapperVersion(file.Content)
	if v == "" {
		return "", fmt.Errorf("Could not determine version from %s", gradleWrapperPropsPath)
	}
	return v, nil
}
