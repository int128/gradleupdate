package cache

import (
	"context"
	"sync"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
)

type ConfigRepository struct {
	ConfigRepository gateways.ConfigRepository
	Logger           gateways.Logger

	l sync.Mutex
	v *domain.Config
}

func (r *ConfigRepository) Get(ctx context.Context) (*domain.Config, error) {
	r.l.Lock()
	defer r.l.Unlock()
	if r.v == nil {
		config, err := r.ConfigRepository.Get(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "error while ConfigRepository.Get")
		}
		r.v = config
		r.Logger.Infof(ctx, "cached Config in memory")
	}
	return r.v, nil
}
