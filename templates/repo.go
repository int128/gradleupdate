package templates

import (
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
)

type Repository struct {
	Repository                  git.Repository
	GradleUpdatePreconditionOut gradle.UpdatePreconditionOut
	BadgeMarkdown               string
	BadgeHTML                   string
	BadgeURL                    string
	RequestUpdateURL            string
}
