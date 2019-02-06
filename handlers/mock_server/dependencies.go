package main

import (
	"context"
	"log"
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
	"go.uber.org/zap"
)

func newContainer() (*dig.Container, error) {
	c := dig.New()
	for _, dependency := range dependencies {
		if err := c.Provide(dependency); err != nil {
			return nil, errors.Wrapf(err, "error while providing dependencies")
		}
	}
	return c, nil
}

type app struct {
	dig.In
	Router handlers.Router
	Logger gateways.Logger
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
				return &usecases.GetBadgeResponse{
					CurrentVersion: gradle.Version("5.0"),
					UpToDate:       false,
				}, nil
			})
		return getBadge
	},
	func(ctrl *gomock.Controller) usecases.GetRepository {
		getRepository := usecaseTestDoubles.NewMockGetRepository(ctrl)
		repository := git.Repository{
			ID:          git.RepositoryID{Owner: "int128", Name: "gradleupdate"},
			Description: "Automatic Gradle Update Service",
			AvatarURL:   "https://avatars0.githubusercontent.com/u/321266",
			HTMLURL:     "https://github.com/int128/gradleupdate",
		}
		getRepository.EXPECT().
			Do(gomock.Not(nil), gomock.Any()).
			AnyTimes().
			DoAndReturn(func(ctx context.Context, id git.RepositoryID) (*usecases.GetRepositoryResponse, error) {
				return &usecases.GetRepositoryResponse{
					Repository:                  repository,
					LatestGradleRelease:         gradle.Release{Version: "5.1"},
					UpdatePreconditionViolation: gradleupdate.ReadyToUpdate,
				}, nil
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

	func() *gomock.Controller {
		return gomock.NewController(&testReporter{})
	},
	func() (gateways.Logger, error) {
		// skip 1st caller that is zapLogger
		logger, err := zap.NewDevelopment(zap.AddCallerSkip(1))
		if err != nil {
			return nil, errors.Wrapf(err, "error while creating a logger")
		}
		return &zapLogger{logger.Sugar()}, nil
	},
}

type zapLogger struct {
	sugar *zap.SugaredLogger
}

func (l *zapLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

func (l *zapLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

func (l *zapLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

func (l *zapLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.sugar.Errorf(format, args)
}

type testReporter struct{}

func (t *testReporter) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (t *testReporter) Fatalf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
