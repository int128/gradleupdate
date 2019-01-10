package templates

import "github.com/int128/gradleupdate/domain"

type Repository struct {
	Repository    domain.Repository
	LatestVersion domain.GradleVersion
	UpToDate      bool

	ThisURL            string
	BadgeURL           string
	SendPullRequestURL string
	BaseURL            string
}
