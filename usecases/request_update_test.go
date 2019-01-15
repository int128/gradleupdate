package usecases_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/testdata"
	"github.com/int128/gradleupdate/gateways/interfaces/mock_gateways"
	"github.com/int128/gradleupdate/usecases"
	interfaces "github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces/mock_usecases"
)

func TestRequestUpdate_Do_UpToDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	repositoryID := domain.RepositoryID{Owner: "owner", Name: "repo"}

	repositoryRepository := mock_gateways.NewMockRepositoryRepository(ctrl)
	repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
		Return(testdata.GradleWrapperProperties4102, nil)

	gradleService := mock_gateways.NewMockGradleService(ctrl)
	gradleService.EXPECT().GetCurrentVersion(ctx).
		Return(domain.GradleVersion("4.10.2"), nil)

	sendPullRequest := mock_usecases.NewMockSendPullRequest(ctrl)

	u := usecases.RequestUpdate{
		RepositoryRepository: repositoryRepository,
		GradleService:        gradleService,
		SendPullRequest:      sendPullRequest,
	}
	err := u.Do(ctx, repositoryID)
	if err != nil {
		t.Fatalf("could not do usecase: %s", err)
	}
}

func TestRequestUpdate_Do_OutOfDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	repositoryID := domain.RepositoryID{Owner: "owner", Name: "repo"}

	repositoryRepository := mock_gateways.NewMockRepositoryRepository(ctrl)
	repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
		Return(testdata.GradleWrapperProperties4102, nil)

	gradleService := mock_gateways.NewMockGradleService(ctrl)
	gradleService.EXPECT().GetCurrentVersion(ctx).
		Return(domain.GradleVersion("5.0"), nil)

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

	u := usecases.RequestUpdate{
		RepositoryRepository: repositoryRepository,
		GradleService:        gradleService,
		SendPullRequest:      sendPullRequest,
	}
	err := u.Do(ctx, repositoryID)
	if err != nil {
		t.Fatalf("could not do usecase: %s", err)
	}
}
