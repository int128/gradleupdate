package config

import (
	"testing"

	"github.com/int128/gradleupdate/domain/git"
)

func TestToggles_IsEligibleForBatchSendUpdates(t *testing.T) {
	for _, c := range []struct {
		toggles  Toggles
		owner    string
		eligible bool
	}{
		{
			Toggles{},
			"foo",
			true,
		}, {
			Toggles{BatchSendUpdatesOwners: []string{"foo"}},
			"foo",
			true,
		}, {
			Toggles{BatchSendUpdatesOwners: []string{"foo"}},
			"bar",
			false,
		}, {
			Toggles{BatchSendUpdatesOwners: []string{"foo", "bar"}},
			"bar",
			true,
		}, {
			Toggles{BatchSendUpdatesOwners: []string{"foo", "bar"}},
			"baz",
			false,
		},
	} {
		actual := c.toggles.IsEligibleForBatchSendUpdates(git.RepositoryID{Owner: c.owner, Name: "bar"})
		if c.eligible != actual {
			t.Errorf("eligible wants %v but %v", c.eligible, actual)
		}
	}
}
