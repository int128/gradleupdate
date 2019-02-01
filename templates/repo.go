package templates

import (
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradleupdate"
)

type Repository struct {
	Repository                  git.Repository
	UpdatePreconditionViolation gradleupdate.PreconditionViolation
	BadgeMarkdown               string
	BadgeHTML                   string
	BadgeURL                    string
	RequestUpdateURL            string
}
