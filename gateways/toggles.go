package gateways

import (
	"context"
	"strings"
	"sync"

	"github.com/int128/gradleupdate/domain/config"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
	"google.golang.org/appengine/datastore"
)

/*
NewToggles returns an implementation of gateways.Toggle.

You can create the following entity to enable the feature toggles,

	* Kind = Toggle
	* Key (string) = DEFAULT

with the following properties:

	* BatchSendUpdatesOwners (string) = comma separated names (everyone if blank)

*/
func NewToggles() gateways.Toggles {
	return &togglesCache{
		Base: &togglesData{},
	}
}

type togglesCache struct {
	Base gateways.Toggles
	l    sync.Mutex
	v    *config.Toggles
}

func (r *togglesCache) Get(ctx context.Context) (*config.Toggles, error) {
	r.l.Lock()
	defer r.l.Unlock()
	if r.v != nil {
		return r.v, nil
	}
	v, err := r.Base.Get(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting toggles")
	}
	r.v = v
	return r.v, nil
}

type togglesData struct{}

func (r *togglesData) Get(ctx context.Context) (*config.Toggles, error) {
	var e togglesEntity
	k := togglesKey(ctx, "DEFAULT")
	if err := datastore.Get(ctx, k, &e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return &config.Toggles{}, nil
		}
		return nil, errors.Wrapf(err, "error while getting the entity")
	}
	if e.BatchSendUpdatesOwners == "" {
		return &config.Toggles{}, nil
	}
	return &config.Toggles{
		BatchSendUpdatesOwners: e.BatchSendUpdatesOwnersAsArray(),
	}, nil
}

func togglesKey(ctx context.Context, name string) *datastore.Key {
	return datastore.NewKey(ctx, "Toggles", name, 0, nil)
}

type togglesEntity struct {
	BatchSendUpdatesOwners string
}

func (e *togglesEntity) BatchSendUpdatesOwnersAsArray() []string {
	if e.BatchSendUpdatesOwners == "" {
		return nil
	}
	return strings.Split(e.BatchSendUpdatesOwners, ",")
}
