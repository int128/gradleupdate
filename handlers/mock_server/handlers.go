package main

import (
	"context"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/handlers"
	"github.com/int128/gradleupdate/usecases/interfaces"
	usecaseTestDoubles "github.com/int128/gradleupdate/usecases/interfaces/test_doubles"
)

func newHandlers(ctrl *gomock.Controller) handlers.Handlers {
	getBadge := usecaseTestDoubles.NewMockGetBadge(ctrl)
	getRepository := usecaseTestDoubles.NewMockGetRepository(ctrl)
	sendUpdate := usecaseTestDoubles.NewMockSendUpdate(ctrl)

	// badge
	getBadge.EXPECT().
		Do(gomock.Not(nil), gomock.Any()).
		AnyTimes().
		DoAndReturn(func(ctx context.Context, id domain.RepositoryID) (*usecases.GetBadgeResponse, error) {
			if id == (domain.RepositoryID{Owner: "int128", Name: "latest-gradle-wrapper"}) {
				return &usecases.GetBadgeResponse{
					CurrentVersion: domain.GradleVersion("5.1"),
					UpToDate:       true,
				}, nil
			}
			return &usecases.GetBadgeResponse{
				CurrentVersion: domain.GradleVersion("5.0"),
				UpToDate:       false,
			}, nil
		})

	// repository page
	repository := domain.Repository{
		ID:          domain.RepositoryID{Owner: "int128", Name: "gradleupdate"},
		Description: "Automatic Gradle Update Service",
		AvatarURL:   "https://avatars0.githubusercontent.com/u/321266",
		HTMLURL:     "https://github.com/int128/gradleupdate",
	}
	getRepository.EXPECT().
		Do(gomock.Not(nil), gomock.Any()).
		AnyTimes().
		DoAndReturn(func(ctx context.Context, id domain.RepositoryID) (*usecases.GetRepositoryResponse, error) {
			return &usecases.GetRepositoryResponse{
				Repository:                  repository,
				GradleUpdatePreconditionOut: domain.ReadyToUpdate,
			}, nil
		})

	return handlers.Handlers{
		GetBadge:      handlers.GetBadge{GetBadge: getBadge},
		GetRepository: handlers.GetRepository{GetRepository: getRepository},
		SendUpdate:    handlers.SendUpdate{SendUpdate: sendUpdate},
	}
}
