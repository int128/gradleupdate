// Package di provides mock dependencies for frontend development and handlers testing.
//
// You need to provide the following dependencies manually:
//
// * gateways.Logger
//
package di

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
	"github.com/int128/gradleupdate/gateways/interfaces"
	gatewaysTestDoubles "github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/int128/gradleupdate/handlers"
	"github.com/int128/gradleupdate/usecases/interfaces"
	usecaseTestDoubles "github.com/int128/gradleupdate/usecases/interfaces/test_doubles"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func New() (*dig.Container, error) {
	c := dig.New()
	for _, dependency := range dependencies {
		if err := c.Provide(dependency); err != nil {
			return nil, errors.Wrapf(err, "error while providing dependencies")
		}
	}
	return c, nil
}

var dependencies = []interface{}{
	handlers.NewRouter,

	func(ctrl *gomock.Controller) usecases.GetBadge {
		getBadge := usecaseTestDoubles.NewMockGetBadge(ctrl)
		getBadge.EXPECT().
			Do(gomock.Not(nil), gomock.Any()).
			AnyTimes().
			DoAndReturn(func(ctx context.Context, id git.RepositoryID) (*usecases.GetBadgeResponse, error) {
				if id == (git.RepositoryID{Owner: "int128", Name: "latest-gradle-wrapper"}) {
					return &usecases.GetBadgeResponse{
						CurrentVersion: gradle.Version("5.1"),
						UpToDate:       true,
					}, nil
				}
				if id.Owner == "int128" {
					return &usecases.GetBadgeResponse{
						CurrentVersion: gradle.Version("5.0"),
						UpToDate:       false,
					}, nil
				}
				return nil, errors.New("not found")
			})
		return getBadge
	},
	func(ctrl *gomock.Controller) usecases.GetRepository {
		getRepository := usecaseTestDoubles.NewMockGetRepository(ctrl)
		repositoryOf := func(id git.RepositoryID) git.Repository {
			return git.Repository{
				ID:          id,
				Description: "Automatic Gradle Update Service",
				AvatarURL:   "https://avatars0.githubusercontent.com/u/321266",
				HTMLURL:     "https://github.com/int128/gradleupdate",
			}
		}
		getRepository.EXPECT().
			Do(gomock.Not(nil), gomock.Any()).
			AnyTimes().
			DoAndReturn(func(ctx context.Context, id git.RepositoryID) (*usecases.GetRepositoryResponse, error) {
				if id == (git.RepositoryID{Owner: "int128", Name: "latest-gradle-wrapper"}) {
					return &usecases.GetRepositoryResponse{
						Repository:                  repositoryOf(id),
						LatestGradleRelease:         gradle.Release{Version: "5.1"},
						UpdatePreconditionViolation: gradleupdate.AlreadyHasLatestGradle,
					}, nil
				}
				if id.Owner == "int128" {
					return &usecases.GetRepositoryResponse{
						Repository:                  repositoryOf(id),
						LatestGradleRelease:         gradle.Release{Version: "5.1"},
						UpdatePreconditionViolation: gradleupdate.ReadyToUpdate,
					}, nil
				}
				err := usecaseTestDoubles.NewMockGetRepositoryError(ctrl)
				err.EXPECT().
					NoSuchRepository().
					AnyTimes().
					Return(true)
				return nil, err
			})
		return getRepository
	},
	func(ctrl *gomock.Controller) usecases.SendUpdate {
		sendUpdate := usecaseTestDoubles.NewMockSendUpdate(ctrl)
		sendUpdate.EXPECT().
			Do(gomock.Not(nil), gomock.Any()).
			AnyTimes().
			DoAndReturn(func(ctx context.Context, id git.RepositoryID) error {
				time.Sleep(3 * time.Second)
				return nil
			})
		return sendUpdate
	},
	func(ctrl *gomock.Controller) usecases.BatchSendUpdates {
		batchSendUpdates := usecaseTestDoubles.NewMockBatchSendUpdates(ctrl)
		return batchSendUpdates
	},

	func(ctrl *gomock.Controller) gateways.ConfigRepository {
		configRepository := gatewaysTestDoubles.NewMockConfigRepository(ctrl)
		configRepository.EXPECT().
			Get(gomock.Not(nil)).
			AnyTimes().
			Return(&domain.Config{
				GitHubToken: "",
				CSRFKey:     "0123456789abcdef0123456789abcdef",
			}, nil)
		return configRepository
	},

	func(tr testReporter) *gomock.Controller {
		return gomock.NewController(&tr)
	},
}

type testReporter struct {
	dig.In
	Logger gateways.Logger
}

func (t *testReporter) Errorf(format string, args ...interface{}) {
	t.Logger.Errorf(context.Background(), format, args...)
}

func (t *testReporter) Fatalf(format string, args ...interface{}) {
	t.Logger.Errorf(context.Background(), format, args...)
}
