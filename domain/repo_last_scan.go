package domain

import "time"

type RepositoryLastScan struct {
	Repository   RepositoryID
	LastScanTime time.Time

	NoGradleVersionError     bool
	NoReadmeBadgeError       bool
	AlreadyLatestGradleError bool
}
