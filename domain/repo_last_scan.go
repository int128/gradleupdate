package domain

import (
	"time"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradleupdate"
)

type RepositoryLastScan struct {
	Repository            git.RepositoryID
	LastScanTime          time.Time
	PreconditionViolation gradleupdate.PreconditionViolation
}
