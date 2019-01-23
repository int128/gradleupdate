package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/testdata"
	"github.com/int128/gradleupdate/gateways/interfaces/mock_gateways"
	"github.com/int128/gradleupdate/usecases"
	interfaces "github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces/mock_usecases"
	"github.com/pkg/errors"
)

func TestSendUpdate_Do(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2019, 1, 21, 16, 43, 0, 0, time.UTC)
	repositoryID := domain.RepositoryID{Owner: "owner", Name: "repo"}
	badgeURL := "/owner/repo/status.svg"

	t.Run("OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := mock_gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(domain.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"), nil)

		repositoryLastScanRepository := mock_gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:   repositoryID,
			LastScanTime: now,
		})

		gradleService := mock_gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentVersion(ctx).Return(domain.GradleVersion("5.0"), nil)

		sendPullRequest := mock_usecases.NewMockSendPullRequest(ctrl)
		sendPullRequest.EXPECT().Do(ctx, interfaces.SendPullRequestRequest{
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
			NowFunc:                      func() time.Time { return now },
		}
		err := u.Do(ctx, repositoryID, badgeURL)
		if err != nil {
			t.Fatalf("error while Do: %+v", err)
		}
	})

	t.Run("AlreadyHasLatestGradle", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := mock_gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)

		repositoryLastScanRepository := mock_gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:               repositoryID,
			LastScanTime:             now,
			AlreadyLatestGradleError: true,
		})

		gradleService := mock_gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentVersion(ctx).Return(domain.GradleVersion("4.10.2"), nil)
		sendPullRequest := mock_usecases.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			NowFunc:                      func() time.Time { return now },
		}
		err := u.Do(ctx, repositoryID, badgeURL)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(interfaces.SendUpdateError)
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

		repositoryRepository := mock_gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(nil, &noSuchEntityError{})

		repositoryLastScanRepository := mock_gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:           repositoryID,
			LastScanTime:         now,
			NoGradleVersionError: true,
		})

		gradleService := mock_gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentVersion(ctx).Return(domain.GradleVersion("4.10.2"), nil)
		sendPullRequest := mock_usecases.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			NowFunc:                      func() time.Time { return now },
		}
		err := u.Do(ctx, repositoryID, badgeURL)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(interfaces.SendUpdateError)
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

		repositoryRepository := mock_gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(domain.FileContent("INVALID"), nil)

		repositoryLastScanRepository := mock_gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:           repositoryID,
			LastScanTime:         now,
			NoGradleVersionError: true,
		})

		gradleService := mock_gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentVersion(ctx).Return(domain.GradleVersion("4.10.2"), nil)
		sendPullRequest := mock_usecases.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			NowFunc:                      func() time.Time { return now },
		}
		err := u.Do(ctx, repositoryID, badgeURL)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(interfaces.SendUpdateError)
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

		repositoryRepository := mock_gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).Return(nil, &noSuchEntityError{})

		repositoryLastScanRepository := mock_gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:         repositoryID,
			LastScanTime:       now,
			NoReadmeBadgeError: true,
		})

		gradleService := mock_gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentVersion(ctx).Return(domain.GradleVersion("5.0"), nil)
		sendPullRequest := mock_usecases.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			NowFunc:                      func() time.Time { return now },
		}
		err := u.Do(ctx, repositoryID, badgeURL)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(interfaces.SendUpdateError)
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

		repositoryRepository := mock_gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(domain.FileContent("INVALID"), nil)

		repositoryLastScanRepository := mock_gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:         repositoryID,
			LastScanTime:       now,
			NoReadmeBadgeError: true,
		})

		gradleService := mock_gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentVersion(ctx).Return(domain.GradleVersion("5.0"), nil)
		sendPullRequest := mock_usecases.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			NowFunc:                      func() time.Time { return now },
		}
		err := u.Do(ctx, repositoryID, badgeURL)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(interfaces.SendUpdateError)
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
