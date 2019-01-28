package gateways

import (
	"context"
	"os"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"google.golang.org/appengine/datastore"
)

const configKind = "Config"

func configKey(ctx context.Context, name string) *datastore.Key {
	return datastore.NewKey(ctx, configKind, name, 0, nil)
}

type configEntity struct {
	GitHubToken string
}

type ConfigRepository struct {
	dig.In
}

func (r *ConfigRepository) Get(ctx context.Context) (*domain.Config, error) {
	var e configEntity
	k := configKey(ctx, "DEFAULT")
	if err := datastore.Get(ctx, k, &e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, errors.Wrapf(&repositoryError{error: err, noSuchEntity: true}, "entity %s did not exist", k)
		}
		return nil, errors.Wrapf(err, "error while getting the entity")
	}
	return &domain.Config{
		GitHubToken: e.GitHubToken,
	}, nil
}

// ConfigResolver resolves Config by datastore or environment variables.
type ConfigResolver struct {
	dig.In
	ConfigRepository gateways.ConfigRepository
}

func (r *ConfigResolver) Get(ctx context.Context) (*domain.Config, error) {
	config, err := r.ConfigRepository.Get(ctx)
	if err != nil {
		config := &domain.Config{
			GitHubToken: os.Getenv("GITHUB_TOKEN"),
		}
		if config.IsZero() {
			return nil, errors.Wrapf(err, "the Config did not exist and env was not defined")
		}
		return config, nil
	}
	return config, nil
}
