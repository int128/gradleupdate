package main

import (
	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/handlers"
	"github.com/int128/gradleupdate/usecases/interfaces"
	usecaseTestDoubles "github.com/int128/gradleupdate/usecases/interfaces/test_doubles"
)

var nonNil = gomock.Not(gomock.Nil())

func newHandlers(ctrl *gomock.Controller) handlers.Handlers {
	getBadge := usecaseTestDoubles.NewMockGetBadge(ctrl)
	getRepository := usecaseTestDoubles.NewMockGetRepository(ctrl)
	sendUpdate := usecaseTestDoubles.NewMockSendUpdate(ctrl)

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
	sendUpdate.EXPECT().Do(nonNil, exampleRepository, "/int128/gradleupdate/status.svg").AnyTimes().Return(nil)

	var latestRepository = domain.RepositoryID{Owner: "int128", Name: "latest-gradle-wrapper"}
	getBadge.EXPECT().Do(nonNil, latestRepository).AnyTimes().Return(&usecases.GetBadgeResponse{
		CurrentVersion: domain.GradleVersion("5.1"),
		UpToDate:       true,
	}, nil)

	return handlers.Handlers{
		GetBadge: handlers.GetBadge{
			GetBadge: getBadge,
		},
		GetRepository: handlers.GetRepository{
			GetRepository: getRepository,
		},
		SendUpdate: handlers.SendUpdate{
			SendUpdate: sendUpdate,
		},
	}
}