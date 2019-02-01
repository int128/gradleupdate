package domain

import (
	"time"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
)

type RepositoryLastScan struct {
	Repository      git.RepositoryID
	LastScanTime    time.Time
	PreconditionOut gradle.UpdatePreconditionOut
}
