package di

import (
	impl "github.com/int128/gradleupdate/gateways"
	"github.com/int128/gradleupdate/gateways/cache"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func provideGateways(c *dig.Container) error {
	if err := c.Provide(func(i impl.RepositoryRepository) gateways.RepositoryRepository { return &i }); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(i impl.PullRequestRepository) gateways.PullRequestRepository { return &i }); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(i impl.GitService) gateways.GitService { return &i }); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(i impl.BadgeLastAccessRepository) gateways.BadgeLastAccessRepository { return &i }); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(i impl.RepositoryLastScanRepository) gateways.RepositoryLastScanRepository { return &i }); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(i impl.GradleService) gateways.GradleService { return &i }); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(i impl.AEResponseCacheRepository) gateways.ResponseCacheRepository { return &i }); err != nil {
		return errors.WithStack(err)
	}
	if err := c.Provide(func(i impl.AELogger) gateways.Logger { return &i }); err != nil {
		return errors.WithStack(err)
	}

	if err := c.Provide(func(i impl.ConfigRepository, l gateways.Logger) gateways.ConfigRepository {
		return &cache.ConfigRepository{
			Logger: l,
			ConfigRepository: &impl.ConfigResolver{
				ConfigRepository: &i,
			},
		}
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
