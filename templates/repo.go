package templates

import "github.com/int128/gradleupdate/domain"

type Repository struct {
	Repository                  domain.Repository
	GradleUpdatePreconditionOut domain.GradleUpdatePreconditionOut
	BadgeMarkdown               string
	BadgeHTML                   string
	BadgeURL                    string
	RequestUpdateURL            string
}
