package di

import (
	"net/http"

	"github.com/google/go-github/v24/github"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
	"go.uber.org/dig"

	"github.com/int128/gradleupdate/gateways"
	gatewaysInterfaces "github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/handlers"
	"github.com/int128/gradleupdate/infrastructure"
	"github.com/int128/gradleupdate/usecases"
	usecasesInterfaces "github.com/int128/gradleupdate/usecases/interfaces"
)

// New returns a container.
func New() (*dig.Container, error) {
	c := dig.New()
	for _, dependency := range dependencies {
		if err := c.Provide(dependency); err != nil {
			return nil, errors.Wrapf(err, "error while providing dependency")
		}
	}
	return c, nil
}

var dependencies = []interface{}{
	// usecases
	func(i usecases.GetRepository) usecasesInterfaces.GetRepository { return &i },
	func(i usecases.GetBadge) usecasesInterfaces.GetBadge { return &i },
	func(i usecases.SendUpdate) usecasesInterfaces.SendUpdate { return &i },
	func(i usecases.BatchSendUpdates) usecasesInterfaces.BatchSendUpdates { return &i },

	// gateways
	func(i gateways.GetRepositoryQuery) gatewaysInterfaces.GetRepositoryQuery { return &i },
	func(i gateways.SendUpdateQuery) gatewaysInterfaces.SendUpdateQuery { return &i },
	func(i gateways.RepositoryRepository) gatewaysInterfaces.RepositoryRepository { return &i },
	func(i gateways.PullRequestRepository) gatewaysInterfaces.PullRequestRepository { return &i },
	func(i gateways.BadgeLastAccessRepository) gatewaysInterfaces.BadgeLastAccessRepository { return &i },
	func(i gateways.GradleReleaseRepository) gatewaysInterfaces.GradleReleaseRepository { return &i },
	gateways.NewToggles,
	gateways.NewCredentials,
	func(i gateways.Time) gatewaysInterfaces.Time { return &i },
	func(i gateways.HTTPCacheRepository) gatewaysInterfaces.HTTPCacheRepository { return &i },
	func(i gateways.Queue) gatewaysInterfaces.Queue { return &i },
	func(i gateways.AELogger) gatewaysInterfaces.Logger { return &i },

	// infrastructure
	handlers.NewRouter,
	handlers.NewRouteResolver,
	func(factory infrastructure.GitHubClientFactory) *github.Client { return factory.NewV3() },
	func(factory infrastructure.GitHubClientFactory) *githubv4.Client { return factory.NewV4() },
	func(factory infrastructure.HTTPClientFactory) *http.Client { return factory.New() },
}
