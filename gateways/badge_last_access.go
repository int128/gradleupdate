package gateways

import (
	"context"
	"time"

	"github.com/int128/gradleupdate/domain"
	"github.com/pkg/errors"
	"google.golang.org/appengine/datastore"
)

const badgeLastAccessKind = "BadgeLastAccess"

func newBadgeLastAccessKey(ctx context.Context, id domain.RepositoryID) *datastore.Key {
	return datastore.NewKey(ctx, badgeLastAccessKind, id.FullName().String(), 0, nil)
}

type badgeLastAccessEntity struct {
	LastAccessTime time.Time
	CurrentVersion string
	LatestVersion  string
}

//DEPRECATED: TODO: remove after migration
type badgeLastAccessEntityOld struct {
	badgeLastAccessEntity
	TargetVersion string
}

type BadgeLastAccessRepository struct{}

func (r *BadgeLastAccessRepository) Save(ctx context.Context, a domain.BadgeLastAccess) error {
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

func (r *BadgeLastAccessRepository) FindBySince(ctx context.Context, since time.Time) ([]domain.BadgeLastAccess, error) {
	q := datastore.NewQuery(badgeLastAccessKind).
		Filter("LastAccessTime >=", since).
		Order("-LastAccessTime")
	var entities []*badgeLastAccessEntityOld
	keys, err := q.GetAll(ctx, &entities)
	if err != nil {
		return nil, errors.Wrapf(err, "could not find entities")
	}
	ret := make([]domain.BadgeLastAccess, 0)
	for i, e := range entities {
		m := badgeLastAccessEntityToModel(keys[i], *e)
		if m == nil {
			continue
		}
		ret = append(ret, *m)
	}
	return ret, nil
}

func badgeLastAccessEntityToModel(k *datastore.Key, e badgeLastAccessEntityOld) *domain.BadgeLastAccess {
	repositoryID := domain.RepositoryFullName(k.StringID()).Parse()
	if repositoryID == nil {
		return nil
	}
	currentVersion := domain.GradleVersion(e.CurrentVersion)
	//TODO: remove when schema migration is done
	if e.TargetVersion != "" {
		currentVersion = domain.GradleVersion(e.TargetVersion)
	}
	return &domain.BadgeLastAccess{
		Repository:     *repositoryID,
		LastAccessTime: e.LastAccessTime,
		CurrentVersion: currentVersion,
		LatestVersion:  domain.GradleVersion(e.LatestVersion),
	}
}
