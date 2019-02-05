package usecases

import (
	"context"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"golang.org/x/sync/errgroup"
)

type GetBadge struct {
	dig.In
	GradleService             gateways.GradleService
	RepositoryRepository      gateways.RepositoryRepository
	BadgeLastAccessRepository gateways.BadgeLastAccessRepository
	TimeService               gateways.TimeService
	Logger                    gateways.Logger
}

func (usecase *GetBadge) Do(ctx context.Context, id git.RepositoryID) (*usecases.GetBadgeResponse, error) {
	var currentVersion gradle.Version
	var latestRelease *gradle.Release
	var eg errgroup.Group
	eg.Go(func() error {
		gradleWrapperProperties, err := usecase.RepositoryRepository.GetFileContent(ctx, id, gradle.WrapperPropertiesPath)
		if err != nil {
			return errors.Wrapf(err, "error while getting gradle-wrapper.properties in %s", id)
		}
		currentVersion = gradle.FindWrapperVersion(gradleWrapperProperties)
		if currentVersion == "" {
			return errors.Errorf("error while finding Gradle version in the properties of %s", id)
		}
		return nil
	})
	eg.Go(func() error {
		var err error
		latestRelease, err = usecase.GradleService.GetCurrentRelease(ctx)
		if err != nil {
			return errors.Wrapf(err, "error while getting the latest Gradle release")
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := usecase.BadgeLastAccessRepository.Save(ctx, domain.BadgeLastAccess{
		Repository:     id,
		LastAccessTime: usecase.TimeService.Now(),
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
