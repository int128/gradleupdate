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

const repositoryLastScanKind = "RepositoryLastScan"

func newRepositoryLastScanKey(ctx context.Context, id git.RepositoryID) *datastore.Key {
	return datastore.NewKey(ctx, repositoryLastScanKind, id.FullName().String(), 0, nil)
}

type repositoryLastScanEntity struct {
	LastScanTime          time.Time
	PreconditionViolation int
}

type RepositoryLastScanRepository struct {
	dig.In
}

func (r *RepositoryLastScanRepository) Save(ctx context.Context, a domain.RepositoryLastScan) error {
	k := newRepositoryLastScanKey(ctx, a.Repository)
	_, err := datastore.Put(ctx, k, &repositoryLastScanEntity{
		LastScanTime:          a.LastScanTime,
		PreconditionViolation: int(a.PreconditionViolation),
	})
	if err != nil {
		return errors.Wrapf(err, "error while saving an entity")
	}
	return nil
}
