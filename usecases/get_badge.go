package usecases

import (
	"context"
	"time"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/gateways"
	"github.com/pkg/errors"
	"google.golang.org/appengine/log"
)

// Badge represents whether the wrapper is up-to-date or out-of-date.
type Badge struct {
	TargetVersion domain.GradleVersion
	LatestVersion domain.GradleVersion
	UpToDate      bool
}

// GetBadge provides a usecase to get status of Gradle wrapper in a repository.
type GetBadge struct {
	GradleService             gateways.GradleService
	RepositoryRepository      gateways.RepositoryRepository
	BadgeLastAccessRepository gateways.BadgeLastAccessRepository
}

// Do performs the usecase.
func (usecase *GetBadge) Do(ctx context.Context, owner, repo string) (Badge, error) {
	targetRepository := domain.RepositoryIdentifier{Owner: owner, Name: repo}
	targetVersion, err := usecase.getVersion(ctx, targetRepository)
	if err != nil {
		return Badge{}, errors.Wrapf(err, "Could not get Gradle version of repository %s", targetRepository)
	}
	latestVersion, err := usecase.GradleService.GetCurrentVersion(ctx)
	if err != nil {
		return Badge{}, errors.Wrapf(err, "Could not get the latest Gradle version")
	}
	if err := usecase.BadgeLastAccessRepository.Put(ctx, domain.BadgeLastAccess{
		Repository:     targetRepository,
		LastAccessTime: time.Now(),
		TargetVersion:  targetVersion,
		LatestVersion:  latestVersion,
	}); err != nil {
		log.Errorf(ctx, "Could not save badge access")
	}
	return Badge{
		TargetVersion: targetVersion,
		LatestVersion: latestVersion,
		UpToDate:      domain.IsUpToDate(targetVersion, latestVersion),
	}, nil
}

func (usecase *GetBadge) getVersion(ctx context.Context, id domain.RepositoryIdentifier) (domain.GradleVersion, error) {
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
