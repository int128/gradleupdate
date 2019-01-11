package usecases

import (
	"context"
	"time"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"google.golang.org/appengine/log"
)

type GetBadge struct {
	GradleService             gateways.GradleService
	RepositoryRepository      gateways.RepositoryRepository
	BadgeLastAccessRepository gateways.BadgeLastAccessRepository
}

func (usecase *GetBadge) Do(ctx context.Context, id domain.RepositoryID) (*usecases.GetBadgeResponse, error) {
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

	if err := usecase.BadgeLastAccessRepository.Put(ctx, domain.BadgeLastAccess{
		Repository:     id,
		LastAccessTime: time.Now(),
		CurrentVersion: currentVersion,
		LatestVersion:  latestVersion,
	}); err != nil {
		log.Errorf(ctx, "could not save badge access")
	}
	return &usecases.GetBadgeResponse{
		CurrentVersion: currentVersion,
		UpToDate:       currentVersion.GreaterOrEqualThan(latestVersion),
	}, nil
}
