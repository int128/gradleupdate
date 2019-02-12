package gradleupdate

import (
	"time"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
)

// BadgeLastAccess represents a last access to a repository.
type BadgeLastAccess struct {
	Repository     git.RepositoryID
	LastAccessTime time.Time
	CurrentVersion gradle.Version
	LatestVersion  gradle.Version
}
