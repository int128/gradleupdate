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
	timeService := &gateways.TimeService{
		NowValue: time.Date(2019, 1, 21, 16, 43, 0, 0, time.UTC),
	}

	t.Run("SuccessfullyUpdated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(domain.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"), nil)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&domain.GradleRelease{Version: "5.0"}, nil)

		repositoryLastScanRepository := gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:   repositoryID,
			LastScanTime: timeService.NowValue,
		})

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		sendPullRequest.EXPECT().Do(ctx, usecaseInterfaces.SendPullRequestRequest{
			Base:           repositoryID,
			HeadBranchName: "gradle-5.0-owner",
			CommitMessage:  "Gradle 5.0",
			CommitFiles: []domain.File{{
				Path:    domain.GradleWrapperPropertiesPath,
				Content: testdata.GradleWrapperProperties50,
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
		err := u.Do(ctx, repositoryID)
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
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(domain.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"), nil)

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&domain.GradleRelease{Version: "4.10.2"}, nil)

		repositoryLastScanRepository := gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:      repositoryID,
			LastScanTime:    timeService.NowValue,
			PreconditionOut: domain.AlreadyHasLatestGradle,
		})

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			TimeService:                  timeService,
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecaseInterfaces.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != domain.AlreadyHasLatestGradle {
			t.Errorf("PreconditionViolation wants %v but %v", domain.AlreadyHasLatestGradle, preconditionViolation)
		}
	})

	t.Run("NoGradleVersion/NoGradleWrapperProperties", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(domain.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"), nil).MaxTimes(1)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(nil, &noSuchEntityError{})

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&domain.GradleRelease{Version: "5.0"}, nil).MaxTimes(1)

		repositoryLastScanRepository := gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:      repositoryID,
			LastScanTime:    timeService.NowValue,
			PreconditionOut: domain.NoGradleWrapperProperties,
		})

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			TimeService:                  timeService,
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecaseInterfaces.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != domain.NoGradleWrapperProperties {
			t.Errorf("PreconditionViolation wants %v but %v", domain.NoGradleWrapperProperties, preconditionViolation)
		}
	})

	t.Run("NoGradleVersion/NoGradleVersion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(domain.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"), nil)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(domain.FileContent("INVALID"), nil)

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&domain.GradleRelease{Version: "5.0"}, nil)

		repositoryLastScanRepository := gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:      repositoryID,
			LastScanTime:    timeService.NowValue,
			PreconditionOut: domain.NoGradleVersion,
		})

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			TimeService:                  timeService,
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecaseInterfaces.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != domain.NoGradleVersion {
			t.Errorf("PreconditionViolation wants %v but %v", domain.NoGradleVersion, preconditionViolation)
		}
	})

	t.Run("NoReadmeBadge/NoReadme", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(nil, &noSuchEntityError{})
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil).MaxTimes(1)

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&domain.GradleRelease{Version: "5.0"}, nil).MaxTimes(1)

		repositoryLastScanRepository := gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:      repositoryID,
			LastScanTime:    timeService.NowValue,
			PreconditionOut: domain.NoReadme,
		})

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			TimeService:                  timeService,
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecaseInterfaces.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != domain.NoReadme {
			t.Errorf("PreconditionViolation wants %v but %v", domain.NoReadme, preconditionViolation)
		}
	})

	t.Run("NoReadmeBadge/NoReadmeBadge", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(domain.FileContent("INVALID"), nil)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&domain.GradleRelease{Version: "5.0"}, nil)

		repositoryLastScanRepository := gateways.NewMockRepositoryLastScanRepository(ctrl)
		repositoryLastScanRepository.EXPECT().Save(ctx, domain.RepositoryLastScan{
			Repository:      repositoryID,
			LastScanTime:    timeService.NowValue,
			PreconditionOut: domain.NoReadmeBadge,
		})

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:         repositoryRepository,
			RepositoryLastScanRepository: repositoryLastScanRepository,
			GradleService:                gradleService,
			SendPullRequest:              sendPullRequest,
			TimeService:                  timeService,
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecaseInterfaces.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != domain.NoReadmeBadge {
			t.Errorf("PreconditionViolation wants %v but %v", domain.NoReadmeBadge, preconditionViolation)
		}
	})
}

type noSuchEntityError struct{}

func (err *noSuchEntityError) Error() string      { return "404" }
func (err *noSuchEntityError) NoSuchEntity() bool { return true }
