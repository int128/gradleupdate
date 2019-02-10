package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/config"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/int128/gradleupdate/usecases"
	"github.com/int128/gradleupdate/usecases/interfaces/test_doubles"
)

func TestBatchSendUpdates_Do(t *testing.T) {
	ctx := context.Background()
	latestGradleRelease := gradle.Release{Version: "5.0"}
	fixedTime := &gatewaysTestDoubles.FixedTime{
		NowValue: time.Date(2019, 1, 21, 16, 43, 0, 0, time.UTC),
	}
	oneMonthAgo := time.Date(2018, 12, 22, 16, 43, 0, 0, time.UTC)
	badge1 := domain.BadgeLastAccess{
		Repository:     git.RepositoryID{Owner: "foo", Name: "repo"},
		CurrentVersion: gradle.Version("4.6"),
		LatestVersion:  gradle.Version("5.0"),
		LastAccessTime: time.Date(2019, 1, 1, 12, 34, 0, 0, time.UTC),
	}
	badge2 := domain.BadgeLastAccess{
		Repository:     git.RepositoryID{Owner: "bar", Name: "repo"},
		CurrentVersion: gradle.Version("4.7"),
		LatestVersion:  gradle.Version("5.0"),
		LastAccessTime: time.Date(2019, 1, 2, 12, 34, 0, 0, time.UTC),
	}

	t.Run("ForAll", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		toggles := gatewaysTestDoubles.NewMockToggles(ctrl)
		toggles.EXPECT().
			Get(ctx).
			Return(&config.Toggles{}, nil)

		gradleService := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleService.EXPECT().
			GetCurrent(ctx).
			Return(&latestGradleRelease, nil)

		badgeLastAccessRepository := gatewaysTestDoubles.NewMockBadgeLastAccessRepository(ctrl)
		badgeLastAccessRepository.EXPECT().
			FindBySince(ctx, oneMonthAgo).
			Return([]domain.BadgeLastAccess{badge1, badge2}, nil)

		sendUpdate := usecasesTestDoubles.NewMockSendUpdate(ctrl)
		sendUpdate.EXPECT().
			Do(ctx, git.RepositoryID{Owner: "foo", Name: "repo"}).
			Return(nil)
		sendUpdate.EXPECT().
			Do(ctx, git.RepositoryID{Owner: "bar", Name: "repo"}).
			Return(nil)

		u := usecases.BatchSendUpdates{
			GradleReleaseRepository:   gradleService,
			BadgeLastAccessRepository: badgeLastAccessRepository,
			SendUpdate:                sendUpdate,
			Toggles:                   toggles,
			Time:                      fixedTime,
			Logger:                    gatewaysTestDoubles.NewLogger(t),
		}
		if err := u.Do(ctx); err != nil {
			t.Fatalf("error while executing usecase: %s", err)
		}
	})

	t.Run("LimitOwnersByFeatureToggle", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		toggles := gatewaysTestDoubles.NewMockToggles(ctrl)
		toggles.EXPECT().
			Get(ctx).
			Return(&config.Toggles{BatchSendUpdatesOwners: []string{"foo"}}, nil)

		gradleService := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleService.EXPECT().
			GetCurrent(ctx).
			Return(&latestGradleRelease, nil)

		badgeLastAccessRepository := gatewaysTestDoubles.NewMockBadgeLastAccessRepository(ctrl)
		badgeLastAccessRepository.EXPECT().
			FindBySince(ctx, oneMonthAgo).
			Return([]domain.BadgeLastAccess{badge1, badge2}, nil)

		sendUpdate := usecasesTestDoubles.NewMockSendUpdate(ctrl)
		sendUpdate.EXPECT().
			Do(ctx, git.RepositoryID{Owner: "foo", Name: "repo"}).
			Return(nil)

		u := usecases.BatchSendUpdates{
			GradleReleaseRepository:   gradleService,
			BadgeLastAccessRepository: badgeLastAccessRepository,
			SendUpdate:                sendUpdate,
			Toggles:                   toggles,
			Time:                      fixedTime,
			Logger:                    gatewaysTestDoubles.NewLogger(t),
		}
		if err := u.Do(ctx); err != nil {
			t.Fatalf("error while executing usecase: %s", err)
		}
	})
}
