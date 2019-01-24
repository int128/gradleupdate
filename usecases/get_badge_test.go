package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/testdata"
	"github.com/int128/gradleupdate/gateways/interfaces/mock_gateways"
	"github.com/int128/gradleupdate/gateways/testing_logger"
	"github.com/int128/gradleupdate/usecases"
)

func TestGetBadge_Do(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	repositoryID := domain.RepositoryID{Owner: "owner", Name: "repo"}
	now := time.Now()

	for _, c := range []struct {
		name           string
		content        domain.FileContent
		currentVersion domain.GradleVersion
		latestVersion  domain.GradleVersion
		upToDate       bool
	}{
		{
			name:           "up-to-date",
			content:        testdata.GradleWrapperProperties4102,
			currentVersion: "4.10.2",
			latestVersion:  "4.10.2",
			upToDate:       true,
		},
		{
			name:           "out-of-date",
			content:        testdata.GradleWrapperProperties4102,
			currentVersion: "4.10.2",
			latestVersion:  "5.1",
			upToDate:       false,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			repositoryRepository := mock_gateways.NewMockRepositoryRepository(ctrl)
			repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
				Return(c.content, nil)

			gradleService := mock_gateways.NewMockGradleService(ctrl)
			gradleService.EXPECT().GetCurrentVersion(ctx).
				Return(c.latestVersion, nil)

			badgeLastAccessRepository := mock_gateways.NewMockBadgeLastAccessRepository(ctrl)
			badgeLastAccessRepository.EXPECT().Save(ctx, domain.BadgeLastAccess{
				Repository:     repositoryID,
				CurrentVersion: c.currentVersion,
				LatestVersion:  c.latestVersion,
				LastAccessTime: now,
			}).Return(nil)

			u := usecases.GetBadge{
				RepositoryRepository:      repositoryRepository,
				GradleService:             gradleService,
				BadgeLastAccessRepository: badgeLastAccessRepository,
				TimeProvider:              func() time.Time { return now },
				Logger:                    testing_logger.New(t),
			}
			resp, err := u.Do(ctx, repositoryID)
			if err != nil {
				t.Fatalf("could not do usecase: %s", err)
			}
			if resp.CurrentVersion != c.currentVersion {
				t.Errorf("CurrentVersion wants %s but %s", c.currentVersion, resp.CurrentVersion)
			}
			if resp.UpToDate != c.upToDate {
				t.Errorf("UpToDate wants %v but %v", c.upToDate, resp.UpToDate)
			}
		})
	}
}
