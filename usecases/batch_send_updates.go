package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type BatchSendUpdates struct {
	dig.In
	GradleService             gateways.GradleService
	BadgeLastAccessRepository gateways.BadgeLastAccessRepository
	SendUpdate                usecases.SendUpdate
	TimeService               gateways.TimeService
	Logger                    gateways.Logger
}

func (usecase *BatchSendUpdates) Do(ctx context.Context) error {
	oneMonthAgo := usecase.TimeService.Now().Add(-1 * 30 * 24 * time.Hour)
	badges, err := usecase.BadgeLastAccessRepository.FindBySince(ctx, oneMonthAgo)
	if err != nil {
		return errors.Wrapf(err, "could not find badges since %s", oneMonthAgo)
	}

	latestVersion, err := usecase.GradleService.GetCurrentVersion(ctx)
	if err != nil {
		return errors.Wrapf(err, "could not get the latest Gradle version")
	}
	for _, badge := range badges {
		if badge.CurrentVersion.GreaterOrEqualThan(latestVersion) {
			usecase.Logger.Infof(ctx, "skip the repository %s because it has the latest Gradle", badge.Repository)
			continue
		}
		//TODO: externalize URL provider
		badgeURL := fmt.Sprintf("/%s/%s/status.svg", badge.Repository.Owner, badge.Repository.Name)
		if err := usecase.SendUpdate.Do(ctx, badge.Repository, badgeURL); err != nil {
			usecase.Logger.Warnf(ctx, "could not send an update for repository %s: %+v", badge.Repository, err)
		}
	}
	return nil
}
