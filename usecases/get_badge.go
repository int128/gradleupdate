package usecases

import (
	"context"
	"fmt"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"golang.org/x/sync/errgroup"
)

type GetBadge struct {
	dig.In
	GradleReleaseRepository   gateways.GradleReleaseRepository
	RepositoryRepository      gateways.RepositoryRepository
	BadgeLastAccessRepository gateways.BadgeLastAccessRepository
	Time                      gateways.Time
	Logger                    gateways.Logger
}

func (usecase *GetBadge) Do(ctx context.Context, id git.RepositoryID) (*usecases.GetBadgeResponse, error) {
	var currentVersion gradle.Version
	var latestRelease *gradle.Release
	var eg errgroup.Group
	eg.Go(func() error {
		gradleWrapperProperties, err := usecase.RepositoryRepository.GetFileContent(ctx, id, gradle.WrapperPropertiesPath)
		if err != nil {
			if usecase.RepositoryRepository.IsNoSuchEntityError(err) {
				return errors.Wrapf(&noGradleVersionError{err}, "no gradle-wrapper.properties in %s", id)
			}
			return errors.Wrapf(err, "error while getting gradle-wrapper.properties in %s", id)
		}
		currentVersion = gradle.FindWrapperVersion(gradleWrapperProperties)
		if currentVersion == "" {
			return errors.WithStack(&noGradleVersionError{fmt.Errorf("no Gradle version in gradle-wrapper.properties of %s", id)})
		}
		return nil
	})
	eg.Go(func() error {
		var err error
		latestRelease, err = usecase.GradleReleaseRepository.GetCurrent(ctx)
		if err != nil {
			return errors.Wrapf(err, "error while getting the latest Gradle release")
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := usecase.BadgeLastAccessRepository.Save(ctx, gradleupdate.BadgeLastAccess{
		Repository:     id,
		LastAccessTime: usecase.Time.Now(),
		CurrentVersion: currentVersion,
		LatestVersion:  latestRelease.Version,
	}); err != nil {
		usecase.Logger.Errorf(ctx, "error while saving badge access")
	}
	return &usecases.GetBadgeResponse{
		CurrentVersion: currentVersion,
		UpToDate:       currentVersion.GreaterOrEqualThan(latestRelease.Version),
	}, nil
}

func (usecase *GetBadge) IsNoGradleVersionError(err error) bool {
	_, ok := errors.Cause(err).(*noGradleVersionError)
	return ok
}

type noGradleVersionError struct {
	error
}
