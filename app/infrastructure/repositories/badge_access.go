package repositories

import (
	"context"
	"github.com/int128/gradleupdate/app/domain"
	"github.com/pkg/errors"
	"google.golang.org/appengine/datastore"
	"time"
)

const badgeLastAccessKind = "BadgeLastAccess"

type badgeLastAccessEntity struct {
	LastAccessTime time.Time
	TargetVersion  string
	LatestVersion  string
}

type BadgeLastAccess struct{}

func (r *BadgeLastAccess) Get(ctx context.Context, id domain.RepositoryIdentifier) (domain.BadgeLastAccess, error) {
	k := datastore.NewKey(ctx, badgeLastAccessKind, id.FullName(), 0, nil)
	var e badgeLastAccessEntity
	err := datastore.Get(ctx, k, &e)
	if err != nil {
		return domain.BadgeLastAccess{}, errors.Wrapf(err, "Could not get the entity")
	}
	return domain.BadgeLastAccess{
		Repository:     id,
		LastAccessTime: e.LastAccessTime,
		TargetVersion:  domain.GradleVersion(e.TargetVersion),
		LatestVersion:  domain.GradleVersion(e.LatestVersion),
	}, nil
}

func (r *BadgeLastAccess) Put(ctx context.Context, a domain.BadgeLastAccess) error {
	k := datastore.NewKey(ctx, badgeLastAccessKind, a.Repository.FullName(), 0, nil)
	_, err := datastore.Put(ctx, k, &badgeLastAccessEntity{
		LastAccessTime: a.LastAccessTime,
		TargetVersion:  string(a.TargetVersion),
		LatestVersion:  string(a.LatestVersion),
	})
	if err != nil {
		return errors.Wrapf(err, "Could not put an entity")
	}
	return nil
}
