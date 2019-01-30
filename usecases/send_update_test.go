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
	usecaseInterfaces "github.com/int128/gradleupdate/usecases/interfaces"
	usecaseTestDoubles "github.com/int128/gradleupdate/usecases/interfaces/test_doubles"
	"github.com/pkg/errors"
)

func TestSendUpdate_Do(t *testing.T) {
	ctx := context.Background()
	repositoryID := domain.RepositoryID{Owner: "owner", Name: "repo"}
	badgeURL := "/owner/repo/status.svg"
	timeService := &gateways.TimeService{
		NowValue: time.Date(2019, 1, 21, 16, 43, 0, 0, time.UTC),
	}

	t.Run("OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(domain.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"), nil)

		repositoryLastScanRepository := gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:   repositoryID,
			LastScanTime: timeService.NowValue,
		})

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&domain.GradleRelease{Version: domain.GradleVersion("5.0")}, nil)

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		sendPullRequest.EXPECT().Do(ctx, usecaseInterfaces.SendPullRequestRequest{
			Base:           repositoryID,
			HeadBranchName: "gradle-5.0-owner",
			CommitMessage:  "Gradle 5.0",
			CommitFiles: []domain.File{{
				Path:    domain.GradleWrapperPropertiesPath,
				Content: domain.FileContent(testdata.GradleWrapperProperties50),
			}},
			Title: "Gradle 5.0",
			Body:  "Gradle 5.0 is available.",
		}).Return(nil)

		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			TimeService:                  timeService,
		}
		err := u.Do(ctx, repositoryID, badgeURL)
		if err != nil {
			t.Fatalf("error while Do: %+v", err)
		}
	})

	t.Run("AlreadyHasLatestGradle", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)

		repositoryLastScanRepository := gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:               repositoryID,
			LastScanTime:             timeService.NowValue,
			AlreadyLatestGradleError: true,
		})

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&domain.GradleRelease{Version: domain.GradleVersion("4.10.2")}, nil)

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			TimeService:                  timeService,
		}
		err := u.Do(ctx, repositoryID, badgeURL)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecaseInterfaces.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		if sendUpdateError.AlreadyHasLatestGradle() != true {
			t.Errorf("AlreadyHasLatestGradle wants true but false")
		}
	})

	t.Run("NoGradleVersion/GradleWrapperPropertiesNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(nil, &noSuchEntityError{})

		repositoryLastScanRepository := gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:           repositoryID,
			LastScanTime:         timeService.NowValue,
			NoGradleVersionError: true,
		})

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&domain.GradleRelease{Version: domain.GradleVersion("4.10.2")}, nil)

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			TimeService:                  timeService,
		}
		err := u.Do(ctx, repositoryID, badgeURL)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecaseInterfaces.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		if sendUpdateError.NoGradleVersion() != true {
			t.Errorf("NoGradleVersion wants true but false")
		}
	})

	t.Run("NoGradleVersion/InvalidGradleWrapperProperties", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(domain.FileContent("INVALID"), nil)

		repositoryLastScanRepository := gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:           repositoryID,
			LastScanTime:         timeService.NowValue,
			NoGradleVersionError: true,
		})

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&domain.GradleRelease{Version: domain.GradleVersion("4.10.2")}, nil)

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			TimeService:                  timeService,
		}
		err := u.Do(ctx, repositoryID, badgeURL)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecaseInterfaces.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		if sendUpdateError.NoGradleVersion() != true {
			t.Errorf("NoGradleVersion wants true but false")
		}
	})

	t.Run("NoReadmeBadge/ReadmeNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).Return(nil, &noSuchEntityError{})

		repositoryLastScanRepository := gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:         repositoryID,
			LastScanTime:       timeService.NowValue,
			NoReadmeBadgeError: true,
		})

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&domain.GradleRelease{Version: domain.GradleVersion("5.0")}, nil)

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			TimeService:                  timeService,
		}
		err := u.Do(ctx, repositoryID, badgeURL)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecaseInterfaces.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		if sendUpdateError.NoReadmeBadge() != true {
			t.Errorf("NoReadmeBadge wants true but false")
		}
	})

	t.Run("NoReadmeBadge/InvalidReadme", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(domain.FileContent("INVALID"), nil)

		repositoryLastScanRepository := gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:         repositoryID,
			LastScanTime:       timeService.NowValue,
			NoReadmeBadgeError: true,
		})

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&domain.GradleRelease{Version: domain.GradleVersion("5.0")}, nil)

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			TimeService:                  timeService,
		}
		err := u.Do(ctx, repositoryID, badgeURL)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecaseInterfaces.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		if sendUpdateError.NoReadmeBadge() != true {
			t.Errorf("NoReadmeBadge wants true but false")
		}
	})
}

type noSuchEntityError struct{}

func (err *noSuchEntityError) Error() string      { return "404" }
func (err *noSuchEntityError) NoSuchEntity() bool { return true }
