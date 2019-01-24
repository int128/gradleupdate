package di

import (
	impl "github.com/int128/gradleupdate/infrastructure"
	"github.com/int128/gradleupdate/infrastructure/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func provideInfrastructure(c *dig.Container) error {
	if err := c.Provide(func(i impl.GitHubClientFactory) infrastructure.GitHubClientFactory { return &i }); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(i impl.GradleClient) infrastructure.GradleClient { return &i }); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
