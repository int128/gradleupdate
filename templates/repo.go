package templates

import (
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
)

type Repository struct {
	Repository                  git.Repository
	LatestGradleRelease         gradle.Release
	UpdatePreconditionViolation gradleupdate.PreconditionViolation
	BadgeMarkdown               string
	BadgeHTML                   string
	BadgeURL                    string
	RequestUpdateURL            string
}
