package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
	"github.com/int128/gradleupdate/domain/testdata"
	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/pkg/errors"
)

func TestGetBadge_Do(t *testing.T) {
	ctx := context.Background()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}
	fixedTime := &gatewaysTestDoubles.FixedTime{
		NowValue: time.Date(2019, 1, 21, 16, 43, 0, 0, time.UTC),
	}

	for _, c := range []struct {
		name           string
		content        git.FileContent
		currentVersion gradle.Version
		latestVersion  gradle.Version
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repositoryRepository := gatewaysTestDoubles.NewMockRepositoryRepository(ctrl)
			repositoryRepository.EXPECT().
				GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
				Return(c.content, nil)

			gradleService := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
			gradleService.EXPECT().
				GetCurrent(ctx).
				Return(&gradle.Release{Version: c.latestVersion}, nil)

			badgeLastAccessRepository := gatewaysTestDoubles.NewMockBadgeLastAccessRepository(ctrl)
			badgeLastAccessRepository.EXPECT().
				Save(ctx, gradleupdate.BadgeLastAccess{
					Repository:     repositoryID,
					CurrentVersion: c.currentVersion,
					LatestVersion:  c.latestVersion,
					LastAccessTime: fixedTime.NowValue,
				}).Return(nil)

			u := GetBadge{
				RepositoryRepository:      repositoryRepository,
				GradleReleaseRepository:   gradleService,
				BadgeLastAccessRepository: badgeLastAccessRepository,
				Time:                      fixedTime,
				Logger:                    gatewaysTestDoubles.NewLogger(t),
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

	t.Run("NoGradleVersion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		noSuchEntityErr := errors.New("no such entity")

		repositoryRepository := gatewaysTestDoubles.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().
			GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
			Return(nil, noSuchEntityErr)
		repositoryRepository.EXPECT().
			IsNoSuchEntityError(noSuchEntityErr).
			Return(true)

		gradleService := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleService.EXPECT().
			GetCurrent(ctx).
			MaxTimes(1)

		badgeLastAccessRepository := gatewaysTestDoubles.NewMockBadgeLastAccessRepository(ctrl)

		u := GetBadge{
			RepositoryRepository:      repositoryRepository,
			GradleReleaseRepository:   gradleService,
			BadgeLastAccessRepository: badgeLastAccessRepository,
			Time:                      fixedTime,
			Logger:                    gatewaysTestDoubles.NewLogger(t),
		}
		resp, err := u.Do(ctx, repositoryID)
		if resp != nil {
			t.Errorf("resp wants nil but %+v", resp)
		}
		if err == nil {
			t.Fatalf("err wants non-nil but nil")
		}
		if !u.IsNoGradleVersionError(err) {
			t.Errorf("IsNoGradleVersionError wants true but false")
		}
	})
}
