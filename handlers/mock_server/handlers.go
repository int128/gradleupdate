package main

import (
	"context"
	"net/http"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/handlers"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces/mock_usecases"
)

func newHandlers(ctx context.Context, ctrl *gomock.Controller) *handlers.Handlers {
	getBadge := mock_usecases.NewMockGetBadge(ctrl)
	getRepository := mock_usecases.NewMockGetRepository(ctrl)
	sendPullRequest := mock_usecases.NewMockSendPullRequest(ctrl)

	var exampleRepository = domain.RepositoryID{Owner: "int128", Name: "gradleupdate"}
	getBadge.EXPECT().Do(ctx, exampleRepository).AnyTimes().Return(&usecases.GetBadgeResponse{
		CurrentVersion: domain.GradleVersion("5.0"),
		UpToDate:       false,
	}, nil)
	getRepository.EXPECT().Do(ctx, exampleRepository).AnyTimes().Return(&usecases.GetRepositoryResponse{
		Repository: domain.Repository{
			ID:          exampleRepository,
			Description: "Automatic Gradle Update Service",
			AvatarURL:   "https://avatars0.githubusercontent.com/u/321266",
		},
		LatestVersion: "5.0",
		UpToDate:      false,
	}, nil)
	sendPullRequest.EXPECT().Do(ctx, exampleRepository).AnyTimes().Return(nil)

	var latestRepository = domain.RepositoryID{Owner: "int128", Name: "latest-gradle-wrapper"}
	getBadge.EXPECT().Do(ctx, latestRepository).AnyTimes().Return(&usecases.GetBadgeResponse{
		CurrentVersion: domain.GradleVersion("5.1"),
		UpToDate:       true,
	}, nil)

	contextProvider := func(_ *http.Request) context.Context { return ctx }
	return &handlers.Handlers{
		GetBadge: handlers.GetBadge{
			ContextProvider: contextProvider,
			GetBadge:        getBadge,
		},
		GetRepository: handlers.GetRepository{
			ContextProvider: contextProvider,
			GetRepository:   getRepository,
		},
		SendPullRequest: handlers.SendPullRequest{
			ContextProvider: contextProvider,
			SendPullRequest: sendPullRequest,
		},
	}
}
