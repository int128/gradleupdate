package di

import (
	"github.com/int128/gradleupdate/handlers"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type App struct {
	dig.In
	Handlers handlers.Handlers
}

// Invoke runs the function with dependencies.
func Invoke(runner func(App)) error {
	c, err := New()
	if err != nil {
		return errors.Wrapf(err, "error while setting up a container")
	}
	if err := c.Invoke(runner); err != nil {
		return errors.Wrapf(err, "error while invoking the runner")
	}
	return nil
}

// New returns a container.
func New() (*dig.Container, error) {
	c := dig.New()
	if err := provideAll(c, infrastructureDependencies); err != nil {
		return nil, errors.Wrapf(err, "error while providing infrastructure")
	}
	if err := provideAll(c, gatewaysDependencies); err != nil {
		return nil, errors.Wrapf(err, "error while providing gateways")
	}
	if err := provideAll(c, usecasesDependencies); err != nil {
		return nil, errors.Wrapf(err, "error while providing usecases")
	}
	return c, nil
}

func provideAll(c *dig.Container, providers []interface{}) error {
	for _, provider := range providers {
		if err := c.Provide(provider); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
