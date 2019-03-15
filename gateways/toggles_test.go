package gateways

import (
	"testing"

	"github.com/favclip/testerator"
	"github.com/go-test/deep"
	"github.com/int128/gradleupdate/domain/config"
	"google.golang.org/appengine/datastore"
)

func TestNewToggles(t *testing.T) {
	_, ctx, err := testerator.SpinUp()
	if err != nil {
		t.Fatalf("could not spin up appengine context: %s", err)
	}
	defer testerator.SpinDown()

	t.Run("NoEntity", func(t *testing.T) {
		toggles := NewToggles()
		ct, err := toggles.Get(ctx)
		if err != nil {
			t.Fatalf("error while Get: %+v", err)
		}
		if diff := deep.Equal(&config.Toggles{}, ct); diff != nil {
			t.Error(diff)
		}
	})

	t.Run("Empty", func(t *testing.T) {
		toggles := NewToggles()
		k := togglesKey(ctx, "DEFAULT")
		if _, err := datastore.Put(ctx, k, &togglesEntity{
			BatchSendUpdatesOwners: "",
		}); err != nil {
			t.Fatalf("error while putting an entity: %s", err)
		}
		ct, err := toggles.Get(ctx)
		if err != nil {
			t.Fatalf("error while Get: %+v", err)
		}
		want := &config.Toggles{}
		if diff := deep.Equal(want, ct); diff != nil {
			t.Error(diff)
		}
	})

	t.Run("ValidString", func(t *testing.T) {
		toggles := NewToggles()
		k := togglesKey(ctx, "DEFAULT")
		if _, err := datastore.Put(ctx, k, &togglesEntity{
			BatchSendUpdatesOwners: "foo,bar",
		}); err != nil {
			t.Fatalf("error while putting an entity: %s", err)
		}
		ct, err := toggles.Get(ctx)
		if err != nil {
			t.Fatalf("error while Get: %+v", err)
		}
		want := &config.Toggles{
			BatchSendUpdatesOwners: []string{"foo", "bar"},
		}
		if diff := deep.Equal(want, ct); diff != nil {
			t.Error(diff)
		}
	})
}
