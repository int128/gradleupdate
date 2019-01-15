package main

import (
	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/handlers"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces/mock_usecases"
)

var nonNil = gomock.Not(gomock.Nil())

func newHandlers(ctrl *gomock.Controller) *handlers.Handlers {
	getBadge := mock_usecases.NewMockGetBadge(ctrl)
	getRepository := mock_usecases.NewMockGetRepository(ctrl)
	requestUpdate := mock_usecases.NewMockRequestUpdate(ctrl)

	var exampleRepository = domain.RepositoryID{Owner: "int128", Name: "gradleupdate"}
	getBadge.EXPECT().Do(nonNil, exampleRepository).AnyTimes().Return(&usecases.GetBadgeResponse{
		CurrentVersion: domain.GradleVersion("5.0"),
		UpToDate:       false,
	}, nil)
	getRepository.EXPECT().Do(nonNil, exampleRepository).AnyTimes().Return(&usecases.GetRepositoryResponse{
		Repository: domain.Repository{
			ID:          exampleRepository,
			Description: "Automatic Gradle Update Service",
			AvatarURL:   "https://avatars0.githubusercontent.com/u/321266",
		},
		LatestVersion:  "5.1",
		CurrentVersion: "5.0",
		UpToDate:       false,
	}, nil)
	requestUpdate.EXPECT().Do(nonNil, exampleRepository).AnyTimes().Return(nil)

	var latestRepository = domain.RepositoryID{Owner: "int128", Name: "latest-gradle-wrapper"}
	getBadge.EXPECT().Do(nonNil, latestRepository).AnyTimes().Return(&usecases.GetBadgeResponse{
		CurrentVersion: domain.GradleVersion("5.1"),
		UpToDate:       true,
	}, nil)

	return &handlers.Handlers{
		GetBadge: handlers.GetBadge{
			GetBadge: getBadge,
		},
		GetRepository: handlers.GetRepository{
			GetRepository: getRepository,
		},
		RequestUpdate: handlers.RequestUpdate{
			RequestUpdate: requestUpdate,
		},
	}
}
