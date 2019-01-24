package gateways

import (
	"context"
	"time"

	"github.com/int128/gradleupdate/domain"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"google.golang.org/appengine/datastore"
)

const repositoryLastScanKind = "RepositoryLastScan"

func newRepositoryLastScanKey(ctx context.Context, id domain.RepositoryID) *datastore.Key {
	return datastore.NewKey(ctx, repositoryLastScanKind, id.FullName().String(), 0, nil)
}

type repositoryLastScanEntity struct {
	LastScanTime             time.Time
	NoGradleVersionError     bool
	NoReadmeBadgeError       bool
	AlreadyLatestGradleError bool
}

type RepositoryLastScanRepository struct {
	dig.In
}

func (r *RepositoryLastScanRepository) Save(ctx context.Context, a domain.RepositoryLastScan) error {
	k := newRepositoryLastScanKey(ctx, a.Repository)
	_, err := datastore.Put(ctx, k, &repositoryLastScanEntity{
		LastScanTime:             a.LastScanTime,
		NoGradleVersionError:     a.NoGradleVersionError,
		NoReadmeBadgeError:       a.NoReadmeBadgeError,
		AlreadyLatestGradleError: a.AlreadyLatestGradleError,
	})
	if err != nil {
		return errors.Wrapf(err, "could not save an entity")
	}
	return nil
}
