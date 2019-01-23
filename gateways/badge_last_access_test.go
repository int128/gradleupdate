package gateways

import (
	"testing"
	"time"

	"github.com/favclip/testerator"
	"github.com/go-test/deep"
	"github.com/int128/gradleupdate/domain"
	"google.golang.org/appengine/datastore"
)

func TestBadgeLastAccessRepository_Save(t *testing.T) {
	_, ctx, err := testerator.SpinUp()
	if err != nil {
		t.Fatalf("could not spin up appengine context: %s", err)
	}
	defer testerator.SpinDown()
	var r BadgeLastAccessRepository
	baseTime := time.Now()

	if err := r.Save(ctx, domain.BadgeLastAccess{
		Repository:     domain.RepositoryID{Owner: "owner", Name: "repo1"},
		LatestVersion:  "5.0",
		CurrentVersion: "4.1",
		LastAccessTime: baseTime,
	}); err != nil {
		t.Fatalf("could not put the BadgeLastAccess: %s", err)
	}

	var entities []*badgeLastAccessEntity
	keys, err := datastore.NewQuery(badgeLastAccessKind).GetAll(ctx, &entities)
	if err != nil {
		t.Fatalf("could not get all entities: %s", err)
	}
	if len(keys) != 1 {
		t.Errorf("len(keys) wants 1 but %d", len(keys))
	}
	if keys[0].StringID() != "owner/repo1" {
		t.Errorf("ID wants %s but %s", "owner/repo1", keys[0].StringID())
	}
	want := []*badgeLastAccessEntity{
		{
			LatestVersion:  "5.0",
			CurrentVersion: "4.1",
			LastAccessTime: baseTime,
		},
	}
	if diff := deep.Equal(want, entities); diff != nil {
		t.Error(diff)
	}
}

func TestBadgeLastAccessRepository_FindBySince(t *testing.T) {
	_, ctx, err := testerator.SpinUp()
	if err != nil {
		t.Fatalf("could not spin up appengine context: %s", err)
	}
	defer testerator.SpinDown()
	var r BadgeLastAccessRepository
	baseTime := time.Now()
	t.Run("setup", func(t *testing.T) {
		keys := []*datastore.Key{
			newBadgeLastAccessKey(ctx, domain.RepositoryID{Owner: "owner", Name: "repo1"}),
			newBadgeLastAccessKey(ctx, domain.RepositoryID{Owner: "owner", Name: "repo2"}),
		}
		entities := []*badgeLastAccessEntity{
			{LastAccessTime: baseTime.Add(-1 * 24 * time.Hour), CurrentVersion: "4.1", LatestVersion: "5.0"},
			{LastAccessTime: baseTime.Add(-2 * 24 * time.Hour), CurrentVersion: "4.2", LatestVersion: "5.0"},
		}
		if _, err := datastore.PutMulti(ctx, keys, entities); err != nil {
			t.Fatalf("could not save entities: %s", err)
		}
		//TODO: remove after migration
		k := newBadgeLastAccessKey(ctx, domain.RepositoryID{Owner: "owner", Name: "repo3"})
		e := badgeLastAccessEntityOld{
			badgeLastAccessEntity{LastAccessTime: baseTime.Add(-3 * 24 * time.Hour), LatestVersion: "5.0"},
			"4.3",
		}
		if _, err := datastore.Put(ctx, k, &e); err != nil {
			t.Fatalf("could not save the entity: %s", err)
		}
	})

	t.Run("0 day ago", func(t *testing.T) {
		found, err := r.FindBySince(ctx, baseTime)
		if err != nil {
			t.Fatalf("could not find entities: %s", err)
		}
		if diff := deep.Equal([]domain.BadgeLastAccess{}, found); diff != nil {
			t.Error(diff)
		}
	})
	t.Run("1 day ago", func(t *testing.T) {
		found, err := r.FindBySince(ctx, baseTime.Add(-1*24*time.Hour))
		if err != nil {
			t.Fatalf("could not find entities: %s", err)
		}
		want := []domain.BadgeLastAccess{
			{
				Repository:     domain.RepositoryID{Owner: "owner", Name: "repo1"},
				LastAccessTime: baseTime.Add(-1 * 24 * time.Hour),
				CurrentVersion: "4.1",
				LatestVersion:  "5.0",
			},
		}
		if diff := deep.Equal(want, found); diff != nil {
			t.Error(diff)
		}
	})
	t.Run("2 days ago", func(t *testing.T) {
		found, err := r.FindBySince(ctx, baseTime.Add(-2*24*time.Hour))
		if err != nil {
			t.Fatalf("could not find entities: %s", err)
		}
		want := []domain.BadgeLastAccess{
			{
				Repository:     domain.RepositoryID{Owner: "owner", Name: "repo1"},
				LastAccessTime: baseTime.Add(-1 * 24 * time.Hour),
				CurrentVersion: "4.1",
				LatestVersion:  "5.0",
			},
			{
				Repository:     domain.RepositoryID{Owner: "owner", Name: "repo2"},
				LastAccessTime: baseTime.Add(-2 * 24 * time.Hour),
				CurrentVersion: "4.2",
				LatestVersion:  "5.0",
			},
		}
		if diff := deep.Equal(want, found); diff != nil {
			t.Error(diff)
		}
	})
	t.Run("3 days ago", func(t *testing.T) {
		found, err := r.FindBySince(ctx, baseTime.Add(-3*24*time.Hour))
		if err != nil {
			t.Fatalf("could not find entities: %s", err)
		}
		want := []domain.BadgeLastAccess{
			{
				Repository:     domain.RepositoryID{Owner: "owner", Name: "repo1"},
				LastAccessTime: baseTime.Add(-1 * 24 * time.Hour),
				CurrentVersion: "4.1",
				LatestVersion:  "5.0",
			},
			{
				Repository:     domain.RepositoryID{Owner: "owner", Name: "repo2"},
				LastAccessTime: baseTime.Add(-2 * 24 * time.Hour),
				CurrentVersion: "4.2",
				LatestVersion:  "5.0",
			},
			{
				Repository:     domain.RepositoryID{Owner: "owner", Name: "repo3"},
				LastAccessTime: baseTime.Add(-3 * 24 * time.Hour),
				CurrentVersion: "4.3",
				LatestVersion:  "5.0",
			},
		}
		if diff := deep.Equal(want, found); diff != nil {
			t.Error(diff)
		}
	})
}
