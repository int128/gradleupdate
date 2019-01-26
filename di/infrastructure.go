package di

import (
	"net/http"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/infrastructure"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func provideInfrastructure(c *dig.Container) error {
	if err := c.Provide(func(f infrastructure.GitHubClientFactory) *github.Client { return f.New() }); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(f infrastructure.HTTPClientFactory) *http.Client { return f.New() }); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
