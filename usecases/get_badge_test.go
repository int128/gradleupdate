package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/testdata"
	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/int128/gradleupdate/usecases"
)

func TestGetBadge_Do(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	repositoryID := domain.RepositoryID{Owner: "owner", Name: "repo"}
	timeService := &gateways.TimeService{
		NowValue: time.Date(2019, 1, 21, 16, 43, 0, 0, time.UTC),
	}

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
			repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
			repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
				Return(c.content, nil)

			gradleService := gateways.NewMockGradleService(ctrl)
			gradleService.EXPECT().GetCurrentRelease(ctx).
				Return(&domain.GradleRelease{Version: c.latestVersion}, nil)

			badgeLastAccessRepository := gateways.NewMockBadgeLastAccessRepository(ctrl)
			badgeLastAccessRepository.EXPECT().Save(ctx, domain.BadgeLastAccess{
				Repository:     repositoryID,
				CurrentVersion: c.currentVersion,
				LatestVersion:  c.latestVersion,
				LastAccessTime: timeService.NowValue,
			}).Return(nil)

			u := usecases.GetBadge{
				RepositoryRepository:      repositoryRepository,
				GradleService:             gradleService,
				BadgeLastAccessRepository: badgeLastAccessRepository,
				TimeService:               timeService,
				Logger:                    gateways.NewLogger(t),
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
