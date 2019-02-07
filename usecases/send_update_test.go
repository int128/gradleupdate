package usecases_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
	"github.com/int128/gradleupdate/domain/testdata"
	gatewaysInterfaces "github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/int128/gradleupdate/usecases"
	usecaseInterfaces "github.com/int128/gradleupdate/usecases/interfaces"
	usecaseTestDoubles "github.com/int128/gradleupdate/usecases/interfaces/test_doubles"
	"github.com/pkg/errors"
)

func TestSendUpdate_Do(t *testing.T) {
	ctx := context.Background()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}
	timeService := &gateways.TimeService{
		NowValue: time.Date(2019, 1, 21, 16, 43, 0, 0, time.UTC),
	}
	readmeContent := git.FileContent("![Gradle Status](https://gradleupdate.appspot.com/owner/repo/status.svg)")

	t.Run("SuccessfullyUpdated", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(readmeContent, nil)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		repositoryLastUpdateRepository := gateways.NewMockRepositoryLastUpdateRepository(ctrl)
		repositoryLastUpdateRepository.EXPECT().Save(ctx, domain.RepositoryLastUpdate{
			Repository:     repositoryID,
			LastUpdateTime: timeService.NowValue,
		})

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		sendPullRequest.EXPECT().Do(ctx, usecaseInterfaces.SendPullRequestRequest{
			Base:           repositoryID,
			HeadBranchName: "gradle-5.0-owner",
			CommitMessage:  "Gradle 5.0",
			CommitFiles: []git.File{{
				Path:    gradle.WrapperPropertiesPath,
				Content: testdata.GradleWrapperProperties50,
			}},
			Title: "Gradle 5.0",
			Body: `Gradle 5.0 is available.

This is sent by @gradleupdate. See https://gradleupdate.appspot.com/owner/repo/status for more.`,
		}).Return(nil)

		u := usecases.SendUpdate{
			RepositoryRepository:           repositoryRepository,
			RepositoryLastUpdateRepository: repositoryLastUpdateRepository,
			GradleService:                  gradleService,
			SendPullRequest:                sendPullRequest,
			TimeService:                    timeService,
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
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(readmeContent, nil)

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&gradle.Release{Version: "4.10.2"}, nil)

		repositoryLastUpdateRepository := gateways.NewMockRepositoryLastUpdateRepository(ctrl)
		repositoryLastUpdateRepository.EXPECT().Save(ctx, domain.RepositoryLastUpdate{
			Repository:            repositoryID,
			LastUpdateTime:        timeService.NowValue,
			PreconditionViolation: gradleupdate.AlreadyHasLatestGradle,
		})

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:           repositoryRepository,
			RepositoryLastUpdateRepository: repositoryLastUpdateRepository,
			GradleService:                  gradleService,
			SendPullRequest:                sendPullRequest,
			TimeService:                    timeService,
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
		if preconditionViolation != gradleupdate.AlreadyHasLatestGradle {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.AlreadyHasLatestGradle, preconditionViolation)
		}
	})

	t.Run("NoGradleVersion/NoGradleWrapperProperties", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(readmeContent, nil).MaxTimes(1)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
			Return(nil, &noSuchEntityError{})

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil).MaxTimes(1)

		repositoryLastUpdateRepository := gateways.NewMockRepositoryLastUpdateRepository(ctrl)
		repositoryLastUpdateRepository.EXPECT().Save(ctx, domain.RepositoryLastUpdate{
			Repository:            repositoryID,
			LastUpdateTime:        timeService.NowValue,
			PreconditionViolation: gradleupdate.NoGradleWrapperProperties,
		})

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:           repositoryRepository,
			RepositoryLastUpdateRepository: repositoryLastUpdateRepository,
			GradleService:                  gradleService,
			SendPullRequest:                sendPullRequest,
			TimeService:                    timeService,
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
		if preconditionViolation != gradleupdate.NoGradleWrapperProperties {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.NoGradleWrapperProperties, preconditionViolation)
		}
	})

	t.Run("NoGradleVersion/NoGradleVersion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(readmeContent, nil)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
			Return(git.FileContent("INVALID"), nil)

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		repositoryLastUpdateRepository := gateways.NewMockRepositoryLastUpdateRepository(ctrl)
		repositoryLastUpdateRepository.EXPECT().Save(ctx, domain.RepositoryLastUpdate{
			Repository:            repositoryID,
			LastUpdateTime:        timeService.NowValue,
			PreconditionViolation: gradleupdate.NoGradleVersion,
		})

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:           repositoryRepository,
			RepositoryLastUpdateRepository: repositoryLastUpdateRepository,
			GradleService:                  gradleService,
			SendPullRequest:                sendPullRequest,
			TimeService:                    timeService,
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
		if preconditionViolation != gradleupdate.NoGradleVersion {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.NoGradleVersion, preconditionViolation)
		}
	})

	t.Run("NoReadmeBadge/NoReadme", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(nil, &noSuchEntityError{})
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil).MaxTimes(1)

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil).MaxTimes(1)

		repositoryLastUpdateRepository := gateways.NewMockRepositoryLastUpdateRepository(ctrl)
		repositoryLastUpdateRepository.EXPECT().Save(ctx, domain.RepositoryLastUpdate{
			Repository:            repositoryID,
			LastUpdateTime:        timeService.NowValue,
			PreconditionViolation: gradleupdate.NoReadme,
		})

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:           repositoryRepository,
			RepositoryLastUpdateRepository: repositoryLastUpdateRepository,
			GradleService:                  gradleService,
			SendPullRequest:                sendPullRequest,
			TimeService:                    timeService,
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
		if preconditionViolation != gradleupdate.NoReadme {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.NoReadme, preconditionViolation)
		}
	})

	t.Run("NoReadmeBadge/NoReadmeBadge", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
		repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
			Return(git.FileContent("INVALID"), nil)
		repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
			Return(testdata.GradleWrapperProperties4102, nil)

		gradleService := gateways.NewMockGradleService(ctrl)
		gradleService.EXPECT().GetCurrentRelease(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		repositoryLastUpdateRepository := gateways.NewMockRepositoryLastUpdateRepository(ctrl)
		repositoryLastUpdateRepository.EXPECT().Save(ctx, domain.RepositoryLastUpdate{
			Repository:            repositoryID,
			LastUpdateTime:        timeService.NowValue,
			PreconditionViolation: gradleupdate.NoReadmeBadge,
		})

		sendPullRequest := usecaseTestDoubles.NewMockSendPullRequest(ctrl)
		u := usecases.SendUpdate{
			RepositoryRepository:           repositoryRepository,
			RepositoryLastUpdateRepository: repositoryLastUpdateRepository,
			GradleService:                  gradleService,
			SendPullRequest:                sendPullRequest,
			TimeService:                    timeService,
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
		if preconditionViolation != gradleupdate.NoReadmeBadge {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.NoReadmeBadge, preconditionViolation)
		}
	})
}

type noSuchEntityError struct{}

func (err *noSuchEntityError) Error() string       { return "404" }
func (err *noSuchEntityError) NoSuchEntity() bool  { return true }
func (err *noSuchEntityError) AlreadyExists() bool { return false }

var _ gatewaysInterfaces.RepositoryError = &noSuchEntityError{}
