package gateways

import (
	"context"
	"time"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"google.golang.org/appengine/datastore"
)

const repositoryLastUpdateKind = "RepositoryLastUpdate"

func newRepositoryLastUpdateKey(ctx context.Context, id git.RepositoryID) *datastore.Key {
	return datastore.NewKey(ctx, repositoryLastUpdateKind, id.FullName().String(), 0, nil)
}

type repositoryLastUpdateEntity struct {
	LastUpdateTime        time.Time
	PreconditionViolation int
}

type RepositoryLastUpdateRepository struct {
	dig.In
}

func (r *RepositoryLastUpdateRepository) Save(ctx context.Context, a domain.RepositoryLastUpdate) error {
	k := newRepositoryLastUpdateKey(ctx, a.Repository)
	_, err := datastore.Put(ctx, k, &repositoryLastUpdateEntity{
		LastUpdateTime:        a.LastUpdateTime,
		PreconditionViolation: int(a.PreconditionViolation),
	})
	if err != nil {
		return errors.Wrapf(err, "error while saving an entity")
	}
	return nil
}
