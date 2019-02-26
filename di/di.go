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
	func(i usecases.GetRepositoryIn) usecasesInterfaces.GetRepository {
		return &usecases.GetRepository{GetRepositoryIn: i}
	},
	func(i usecases.GetBadge) usecasesInterfaces.GetBadge { return &i },
	func(i usecases.SendUpdateIn) usecasesInterfaces.SendUpdate {
		return &usecases.SendUpdate{SendUpdateIn: i}
	},
	func(i usecases.BatchSendUpdates) usecasesInterfaces.BatchSendUpdates { return &i },

	// gateways
	func(i gateways.GetRepositoryQueryIn) gatewaysInterfaces.GetRepositoryQuery {
		return &gateways.GetRepositoryQuery{GetRepositoryQueryIn: i}
	},
	func(i gateways.SendUpdateQueryIn) gatewaysInterfaces.SendUpdateQuery {
		return &gateways.SendUpdateQuery{SendUpdateQueryIn: i}
	},
	func(i gateways.RepositoryRepositoryIn) gatewaysInterfaces.RepositoryRepository {
		return &gateways.RepositoryRepository{RepositoryRepositoryIn: i}
	},
	func(i gateways.PullRequestRepositoryIn) gatewaysInterfaces.PullRequestRepository {
		return &gateways.PullRequestRepository{PullRequestRepositoryIn: i}
	},
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
