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
	"github.com/int128/gradleupdate/domain/config"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/int128/gradleupdate/handlers"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces/test_doubles"
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
	handlers.NewRouteResolver,

	func(ctrl *gomock.Controller) usecases.GetBadge {
		noGradleVersionErr := errors.New("no Gradle version")
		getBadge := usecasesTestDoubles.NewMockGetBadge(ctrl)
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
				return nil, noGradleVersionErr
			})
		getBadge.EXPECT().
			IsNoGradleVersionError(noGradleVersionErr).
			Return(true).
			AnyTimes()
		return getBadge
	},
	func(ctrl *gomock.Controller) usecases.GetRepository {
		noSuchRepositoryErr := errors.New("no such repository")
		getRepository := usecasesTestDoubles.NewMockGetRepository(ctrl)
		repositoryOf := func(id git.RepositoryID) git.Repository {
			return git.Repository{
				ID:          id,
				Description: "Automatic Gradle Update Service",
				AvatarURL:   "https://avatars0.githubusercontent.com/u/321266",
				URL:         "https://github.com/int128/gradleupdate",
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
				return nil, noSuchRepositoryErr
			})
		getRepository.EXPECT().
			IsNoSuchRepositoryError(noSuchRepositoryErr).
			Return(true).
			AnyTimes()
		return getRepository
	},
	func(ctrl *gomock.Controller) usecases.SendUpdate {
		sendUpdate := usecasesTestDoubles.NewMockSendUpdate(ctrl)
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
		batchSendUpdates := usecasesTestDoubles.NewMockBatchSendUpdates(ctrl)
		return batchSendUpdates
	},

	func(ctrl *gomock.Controller) gateways.Credentials {
		credentials := gatewaysTestDoubles.NewMockCredentials(ctrl)
		credentials.EXPECT().
			Get(gomock.Not(nil)).
			AnyTimes().
			Return(&config.Credentials{
				GitHubToken: "",
				CSRFKey:     []byte("0123456789abcdef0123456789abcdef"),
			}, nil)
		return credentials
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
