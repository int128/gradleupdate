package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/int128/gradleupdate/usecases"
	usecaseTestDoubles "github.com/int128/gradleupdate/usecases/interfaces/test_doubles"
)

func TestBatchSendUpdates_Do(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo1"}

	timeService := &gateways.TimeService{
		NowValue: time.Date(2019, 1, 21, 16, 43, 0, 0, time.UTC),
	}

	gradleService := gateways.NewMockGradleService(ctrl)
	gradleService.EXPECT().GetCurrentRelease(ctx).
		Return(&gradle.Release{Version: "5.0"}, nil)

	oneMonthAgo := time.Date(2018, 12, 22, 16, 43, 0, 0, time.UTC)
	badgeLastAccessRepository := gateways.NewMockBadgeLastAccessRepository(ctrl)
	badgeLastAccessRepository.EXPECT().FindBySince(ctx, oneMonthAgo).Return([]domain.BadgeLastAccess{
		{
			Repository:     repositoryID,
			CurrentVersion: gradle.Version("4.6"),
			LatestVersion:  gradle.Version("5.0"),
			LastAccessTime: time.Date(2019, 1, 1, 12, 34, 0, 0, time.UTC),
		},
	}, nil)

	sendUpdate := usecaseTestDoubles.NewMockSendUpdate(ctrl)
	sendUpdate.EXPECT().Do(ctx, repositoryID).Return(nil)

	u := usecases.BatchSendUpdates{
		GradleService:             gradleService,
		BadgeLastAccessRepository: badgeLastAccessRepository,
		SendUpdate:                sendUpdate,
		TimeService:               timeService,
		Logger:                    gateways.NewLogger(t),
	}
	if err := u.Do(ctx); err != nil {
		t.Fatalf("could not do the use case: %s", err)
	}
}
