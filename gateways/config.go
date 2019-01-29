package gateways

import (
	"context"
	"os"
	"sync"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
	"google.golang.org/appengine/datastore"
)

func NewConfigRepository(logger gateways.Logger) gateways.ConfigRepository {
	return &configResolver{
		Logger: logger,
		Sources: []gateways.ConfigRepository{
			&envConfigRepository{},
			&datastoreConfigRepository{},
		},
	}
}

type configResolver struct {
	Logger  gateways.Logger
	Sources []gateways.ConfigRepository

	l sync.Mutex
	v *domain.Config
}

func (r *configResolver) Get(ctx context.Context) (*domain.Config, error) {
	r.l.Lock()
	defer r.l.Unlock()
	if r.v != nil {
		return r.v, nil
	}
	for _, source := range r.Sources {
		config, err := source.Get(ctx)
		if err != nil {
			r.Logger.Warnf(ctx, "error while getting config from %T: %s", source, err)
			continue
		}
		r.v = config
		return r.v, nil
	}
	return nil, errors.New("could not get config from any source")
}

type envConfigRepository struct{}

func (r *envConfigRepository) Get(ctx context.Context) (*domain.Config, error) {
	config := domain.Config{
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
	}
	if config.IsZero() {
		return nil, errors.New("environment variable was not defined")
	}
	return &config, nil
}

const configKind = "Config"

func configKey(ctx context.Context, name string) *datastore.Key {
	return datastore.NewKey(ctx, configKind, name, 0, nil)
}

type configEntity struct {
	GitHubToken string
}

type datastoreConfigRepository struct{}

func (r *datastoreConfigRepository) Get(ctx context.Context) (*domain.Config, error) {
	var e configEntity
	k := configKey(ctx, "DEFAULT")
	if err := datastore.Get(ctx, k, &e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, errors.Wrapf(&repositoryError{error: err, noSuchEntity: true}, "entity(%s) did not exist", k)
		}
		return nil, errors.Wrapf(err, "error while getting the entity")
	}
	return &domain.Config{
		GitHubToken: e.GitHubToken,
	}, nil
}
