package gateways

import (
	"context"
	"time"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"google.golang.org/appengine/datastore"
)

const badgeLastAccessKind = "BadgeLastAccess"

func newBadgeLastAccessKey(ctx context.Context, id git.RepositoryID) *datastore.Key {
	return datastore.NewKey(ctx, badgeLastAccessKind, id.FullName().String(), 0, nil)
}

type badgeLastAccessEntity struct {
	LastAccessTime time.Time
	CurrentVersion string
	LatestVersion  string
}

type BadgeLastAccessRepository struct {
	dig.In
}

func (r *BadgeLastAccessRepository) Save(ctx context.Context, a gradleupdate.BadgeLastAccess) error {
	k := newBadgeLastAccessKey(ctx, a.Repository)
	_, err := datastore.Put(ctx, k, &badgeLastAccessEntity{
		LastAccessTime: a.LastAccessTime,
		CurrentVersion: string(a.CurrentVersion),
		LatestVersion:  string(a.LatestVersion),
	})
	if err != nil {
		return errors.Wrapf(err, "could not put an entity")
	}
	return nil
}

func (r *BadgeLastAccessRepository) FindBySince(ctx context.Context, since time.Time) ([]gradleupdate.BadgeLastAccess, error) {
	q := datastore.NewQuery(badgeLastAccessKind).
		Filter("LastAccessTime >=", since).
		Order("-LastAccessTime")
	var entities []*badgeLastAccessEntity
	keys, err := q.GetAll(ctx, &entities)
	if err != nil {
		return nil, errors.Wrapf(err, "could not find entities")
	}
	ret := make([]gradleupdate.BadgeLastAccess, 0)
	for i, e := range entities {
		m := badgeLastAccessEntityToModel(keys[i], *e)
		if m == nil {
			continue
		}
		ret = append(ret, *m)
	}
	return ret, nil
}

func badgeLastAccessEntityToModel(k *datastore.Key, e badgeLastAccessEntity) *gradleupdate.BadgeLastAccess {
	repositoryID := git.RepositoryFullName(k.StringID()).Parse()
	if repositoryID == nil {
		return nil
	}
	currentVersion := gradle.Version(e.CurrentVersion)
	return &gradleupdate.BadgeLastAccess{
		Repository:     *repositoryID,
		LastAccessTime: e.LastAccessTime,
		CurrentVersion: currentVersion,
		LatestVersion:  gradle.Version(e.LatestVersion),
	}
}
