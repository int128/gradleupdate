package gateways

import (
	"testing"
	"time"

	"github.com/favclip/testerator"
	"github.com/go-test/deep"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradleupdate"
	"google.golang.org/appengine/datastore"
)

func TestBadgeLastAccessRepository_Save(t *testing.T) {
	_, ctx, err := testerator.SpinUp()
	if err != nil {
		t.Fatalf("could not spin up appengine context: %s", err)
	}
	defer testerator.SpinDown()
	var r BadgeLastAccessRepository
	now := time.Date(2019, 1, 21, 16, 43, 0, 0, time.UTC)

	if err := r.Save(ctx, gradleupdate.BadgeLastAccess{
		Repository:     git.RepositoryID{Owner: "owner", Name: "repo1"},
		LatestVersion:  "5.0",
		CurrentVersion: "4.1",
		LastAccessTime: now,
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
			LastAccessTime: now,
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
	now := time.Date(2019, 1, 21, 16, 43, 0, 0, time.UTC)

	t.Run("setup", func(t *testing.T) {
		keys := []*datastore.Key{
			newBadgeLastAccessKey(ctx, git.RepositoryID{Owner: "owner", Name: "repo1"}),
			newBadgeLastAccessKey(ctx, git.RepositoryID{Owner: "owner", Name: "repo2"}),
			newBadgeLastAccessKey(ctx, git.RepositoryID{Owner: "owner", Name: "repo3"}),
		}
		entities := []*badgeLastAccessEntity{
			{LastAccessTime: now.Add(-1 * 24 * time.Hour), CurrentVersion: "4.1", LatestVersion: "5.0"},
			{LastAccessTime: now.Add(-2 * 24 * time.Hour), CurrentVersion: "4.2", LatestVersion: "5.0"},
			{LastAccessTime: now.Add(-3 * 24 * time.Hour), CurrentVersion: "4.3", LatestVersion: "5.0"},
		}
		if _, err := datastore.PutMulti(ctx, keys, entities); err != nil {
			t.Fatalf("could not save entities: %s", err)
		}
	})

	t.Run("0 day ago", func(t *testing.T) {
		found, err := r.FindBySince(ctx, now)
		if err != nil {
			t.Fatalf("could not find entities: %s", err)
		}
		if diff := deep.Equal([]gradleupdate.BadgeLastAccess{}, found); diff != nil {
			t.Error(diff)
		}
	})
	t.Run("1 day ago", func(t *testing.T) {
		found, err := r.FindBySince(ctx, now.Add(-1*24*time.Hour))
		if err != nil {
			t.Fatalf("could not find entities: %s", err)
		}
		want := []gradleupdate.BadgeLastAccess{
			{
				Repository:     git.RepositoryID{Owner: "owner", Name: "repo1"},
				LastAccessTime: now.Add(-1 * 24 * time.Hour),
				CurrentVersion: "4.1",
				LatestVersion:  "5.0",
			},
		}
		if diff := deep.Equal(want, found); diff != nil {
			t.Error(diff)
		}
	})
	t.Run("2 days ago", func(t *testing.T) {
		found, err := r.FindBySince(ctx, now.Add(-2*24*time.Hour))
		if err != nil {
			t.Fatalf("could not find entities: %s", err)
		}
		want := []gradleupdate.BadgeLastAccess{
			{
				Repository:     git.RepositoryID{Owner: "owner", Name: "repo1"},
				LastAccessTime: now.Add(-1 * 24 * time.Hour),
				CurrentVersion: "4.1",
				LatestVersion:  "5.0",
			},
			{
				Repository:     git.RepositoryID{Owner: "owner", Name: "repo2"},
				LastAccessTime: now.Add(-2 * 24 * time.Hour),
				CurrentVersion: "4.2",
				LatestVersion:  "5.0",
			},
		}
		if diff := deep.Equal(want, found); diff != nil {
			t.Error(diff)
		}
	})
	t.Run("3 days ago", func(t *testing.T) {
		found, err := r.FindBySince(ctx, now.Add(-3*24*time.Hour))
		if err != nil {
			t.Fatalf("could not find entities: %s", err)
		}
		want := []gradleupdate.BadgeLastAccess{
			{
				Repository:     git.RepositoryID{Owner: "owner", Name: "repo1"},
				LastAccessTime: now.Add(-1 * 24 * time.Hour),
				CurrentVersion: "4.1",
				LatestVersion:  "5.0",
			},
			{
				Repository:     git.RepositoryID{Owner: "owner", Name: "repo2"},
				LastAccessTime: now.Add(-2 * 24 * time.Hour),
				CurrentVersion: "4.2",
				LatestVersion:  "5.0",
			},
			{
				Repository:     git.RepositoryID{Owner: "owner", Name: "repo3"},
				LastAccessTime: now.Add(-3 * 24 * time.Hour),
				CurrentVersion: "4.3",
				LatestVersion:  "5.0",
			},
		}
		if diff := deep.Equal(want, found); diff != nil {
			t.Error(diff)
		}
	})
}
