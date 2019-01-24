package di

import (
	impl "github.com/int128/gradleupdate/usecases"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func provideUsecases(c *dig.Container) error {
	if err := c.Provide(func(i impl.GetRepository) usecases.GetRepository { return &i }); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(i impl.GetBadge) usecases.GetBadge { return &i }); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(i impl.SendUpdate) usecases.SendUpdate { return &i }); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(i impl.BatchSendUpdates) usecases.BatchSendUpdates { return &i }); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(i impl.SendPullRequest) usecases.SendPullRequest { return &i }); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
