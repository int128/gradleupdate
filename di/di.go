package di

import (
	"github.com/int128/gradleupdate/handlers"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// New returns a container.
func New() (*dig.Container, error) {
	c := dig.New()
	if err := provideInfrastructure(c); err != nil {
		return nil, errors.Wrapf(err, "error while providing infrastructure")
	}
	if err := provideGateways(c); err != nil {
		return nil, errors.Wrapf(err, "error while providing gateways")
	}
	if err := provideUsecases(c); err != nil {
		return nil, errors.Wrapf(err, "error while providing usecases")
	}
	return c, nil
}

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
