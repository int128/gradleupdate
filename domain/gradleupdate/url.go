package gradleupdate

import (
	"fmt"

	"github.com/int128/gradleupdate/domain/git"
)

type BadgeURL string

// NewBadgeURL returns a BadgeURL for the repository.
// For now hard code domain and path here because it has been public and may not be changed.
func NewBadgeURL(id git.RepositoryID) BadgeURL {
	return BadgeURL(fmt.Sprintf("https://gradleupdate.appspot.com/%s/%s/status.svg", id.Owner, id.Name))
}

type RepositoryURL string

// NewRepositoryURL returns a RepositoryURL for the repository.
// For now hard code domain and path here because it has been public and may not be changed.
func NewRepositoryURL(id git.RepositoryID) RepositoryURL {
	return RepositoryURL(fmt.Sprintf("https://gradleupdate.appspot.com/%s/%s/status", id.Owner, id.Name))
}
