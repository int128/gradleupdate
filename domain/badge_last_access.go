package domain

import "time"

// BadgeLastAccess represents a last access to a repository.
type BadgeLastAccess struct {
	Repository     RepositoryID
	LastAccessTime time.Time
	CurrentVersion GradleVersion
	LatestVersion  GradleVersion
}
