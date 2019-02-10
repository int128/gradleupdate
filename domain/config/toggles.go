package config

import "github.com/int128/gradleupdate/domain/git"

// Toggles represents feature toggles for beta testing.
type Toggles struct {
	BatchSendUpdatesOwners []string // if set, BatchSendUpdates sends for only owners
}

func (f Toggles) IsEligibleForBatchSendUpdates(id git.RepositoryID) bool {
	if len(f.BatchSendUpdatesOwners) == 0 {
		return true
	}
	for _, o := range f.BatchSendUpdatesOwners {
		if o == id.Owner {
			return true
		}
	}
	return false
}
