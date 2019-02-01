package domain

import (
	"time"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradleupdate"
)

type RepositoryLastUpdate struct {
	Repository            git.RepositoryID
	LastUpdateTime        time.Time
	PreconditionViolation gradleupdate.PreconditionViolation
}
