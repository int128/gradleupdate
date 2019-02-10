package usecases

import (
	"context"
	"time"

	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type BatchSendUpdates struct {
	dig.In
	GradleReleaseRepository   gateways.GradleReleaseRepository
	BadgeLastAccessRepository gateways.BadgeLastAccessRepository
	SendUpdate                usecases.SendUpdate
	Time                      gateways.Time
	Toggles                   gateways.Toggles
	Logger                    gateways.Logger
}

func (usecase *BatchSendUpdates) Do(ctx context.Context) error {
	oneMonthAgo := usecase.Time.Now().Add(-1 * 30 * 24 * time.Hour)
	badges, err := usecase.BadgeLastAccessRepository.FindBySince(ctx, oneMonthAgo)
	if err != nil {
		return errors.Wrapf(err, "could not find badges since %s", oneMonthAgo)
	}

	latestRelease, err := usecase.GradleReleaseRepository.GetCurrent(ctx)
	if err != nil {
		return errors.Wrapf(err, "could not get the latest Gradle version")
	}
	toggles, err := usecase.Toggles.Get(ctx)
	if err != nil {
		return errors.Wrapf(err, "error while getting config")
	}
	for _, badge := range badges {
		if badge.CurrentVersion.GreaterOrEqualThan(latestRelease.Version) {
			usecase.Logger.Infof(ctx, "skip the repository %s because it has the latest Gradle", badge.Repository)
			continue
		}
		if !toggles.IsEligibleForBatchSendUpdates(badge.Repository) {
			usecase.Logger.Infof(ctx, "skip the repository %s due to the feature toggle", badge.Repository)
			continue
		}
		if err := usecase.SendUpdate.Do(ctx, badge.Repository); err != nil {
			usecase.Logger.Errorf(ctx, "could not send an update for repository %s: %+v", badge.Repository, err)
		}
	}
	return nil
}
