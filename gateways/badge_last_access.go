package gateways

import (
	"context"
	"time"

	"github.com/int128/gradleupdate/domain"
	"github.com/pkg/errors"
	"google.golang.org/appengine/datastore"
)

const badgeLastAccessKind = "BadgeLastAccess"

type badgeLastAccessEntity struct {
	LastAccessTime time.Time
	TargetVersion  string
	LatestVersion  string
}

type BadgeLastAccessRepository struct{}

func (r *BadgeLastAccessRepository) Get(ctx context.Context, id domain.RepositoryID) (*domain.BadgeLastAccess, error) {
	k := datastore.NewKey(ctx, badgeLastAccessKind, id.FullName(), 0, nil)
	var e badgeLastAccessEntity
	err := datastore.Get(ctx, k, &e)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get the entity")
	}
	return &domain.BadgeLastAccess{
		Repository:     id,
		LastAccessTime: e.LastAccessTime,
		CurrentVersion: domain.GradleVersion(e.TargetVersion),
		LatestVersion:  domain.GradleVersion(e.LatestVersion),
	}, nil
}

func (r *BadgeLastAccessRepository) Put(ctx context.Context, a domain.BadgeLastAccess) error {
	k := datastore.NewKey(ctx, badgeLastAccessKind, a.Repository.FullName(), 0, nil)
	_, err := datastore.Put(ctx, k, &badgeLastAccessEntity{
		LastAccessTime: a.LastAccessTime,
		TargetVersion:  string(a.CurrentVersion),
		LatestVersion:  string(a.LatestVersion),
	})
	if err != nil {
		return errors.Wrapf(err, "Could not put an entity")
	}
	return nil
}
