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
	errs := make([]error, 0)
	for _, source := range r.Sources {
		config, err := source.Get(ctx)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "error while getting config from %T", source))
			continue
		}
		r.v = config
		return r.v, nil
	}
	return nil, errors.Errorf("could not get config from any source: %v", errs)
}

type envConfigRepository struct{}

func (r *envConfigRepository) Get(ctx context.Context) (*domain.Config, error) {
	config := domain.Config{
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
		CSRFKey:     os.Getenv("CSRF_KEY"),
	}
	if !config.IsValid() {
		return nil, errors.New("environment variables are not defined")
	}
	return &config, nil
}

const configKind = "Config"

func configKey(ctx context.Context, name string) *datastore.Key {
	return datastore.NewKey(ctx, configKind, name, 0, nil)
}

type configEntity struct {
	GitHubToken string
	CSRFKey     string
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
	config := domain.Config{
		GitHubToken: e.GitHubToken,
		CSRFKey:     e.CSRFKey,
	}
	if !config.IsValid() {
		return nil, errors.New("datastore Config are not defined")
	}
	return &config, nil
}
