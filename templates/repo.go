package templates

import "github.com/int128/gradleupdate/domain"

type Repository struct {
	Repository     domain.Repository
	CurrentVersion domain.GradleVersion
	LatestVersion  domain.GradleVersion
	UpToDate       bool

	ThisURL          string
	BadgeURL         string
	RequestUpdateURL string
	BaseURL          string
}
