package di

import (
	"os"

	"github.com/int128/gradleupdate/gateways/interfaces"
	impl "github.com/int128/gradleupdate/infrastructure"
	"github.com/int128/gradleupdate/infrastructure/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func provideInfrastructure(c *dig.Container) error {
	if err := c.Provide(func(responseCacheRepository gateways.ResponseCacheRepository) infrastructure.GitHubClientFactory {
		return &impl.GitHubClientFactory{
			Token:                   os.Getenv("GITHUB_TOKEN"),
			ResponseCacheRepository: responseCacheRepository,
		}
	}); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(i impl.GradleClient) infrastructure.GradleClient { return &i }); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
