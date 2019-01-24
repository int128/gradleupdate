package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces/mock_gateways"
	"github.com/int128/gradleupdate/gateways/testing_logger"
	"github.com/int128/gradleupdate/usecases"
	"github.com/int128/gradleupdate/usecases/interfaces/mock_usecases"
)

func TestBatchSendUpdates_Do(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	now := time.Date(2019, 1, 21, 16, 43, 0, 0, time.UTC)
	oneMonthAgo := time.Date(2018, 12, 22, 16, 43, 0, 0, time.UTC)
	repositoryID := domain.RepositoryID{Owner: "owner", Name: "repo1"}

	gradleService := mock_gateways.NewMockGradleService(ctrl)
	gradleService.EXPECT().GetCurrentVersion(ctx).
		Return(domain.GradleVersion("5.0"), nil)

	badgeLastAccessRepository := mock_gateways.NewMockBadgeLastAccessRepository(ctrl)
	badgeLastAccessRepository.EXPECT().FindBySince(ctx, oneMonthAgo).Return([]domain.BadgeLastAccess{
		{
			Repository:     repositoryID,
			CurrentVersion: domain.GradleVersion("4.6"),
			LatestVersion:  domain.GradleVersion("5.0"),
			LastAccessTime: time.Date(2019, 1, 1, 12, 34, 0, 0, time.UTC),
		},
	}, nil)

	sendUpdate := mock_usecases.NewMockSendUpdate(ctrl)
	sendUpdate.EXPECT().Do(ctx, repositoryID, "/owner/repo1/status.svg").Return(nil)

	u := usecases.BatchSendUpdates{
		GradleService:             gradleService,
		BadgeLastAccessRepository: badgeLastAccessRepository,
		SendUpdate:                sendUpdate,
		TimeProvider:              func() time.Time { return now },
		Logger:                    testing_logger.New(t),
	}
	if err := u.Do(ctx); err != nil {
		t.Fatalf("could not do the use case: %s", err)
	}
}
