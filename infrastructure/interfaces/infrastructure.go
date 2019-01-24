package infrastructure

import (
	"context"

	"github.com/google/go-github/v18/github"
)

type GitHubClientFactory interface {
	New(ctx context.Context) *github.Client
}

type GradleClient interface {
	GetCurrentVersion(ctx context.Context) (*CurrentVersionResponse, error)
}

type CurrentVersionResponse struct {
	Version            string `json:"version"`
	BuildTime          string `json:"buildTime"`
	Current            bool   `json:"current"`
	Snapshot           bool   `json:"snapshot"`
	Nightly            bool   `json:"nightly"`
	ReleaseNightly     bool   `json:"releaseNightly"`
	ActiveRc           bool   `json:"activeRc"`
	RcFor              string `json:"rcFor"`
	MilestoneFor       string `json:"milestoneFor"`
	Broken             bool   `json:"broken"`
	DownloadURL        string `json:"downloadUrl"`
	ChecksumURL        string `json:"checksumUrl"`
	WrapperChecksumURL string `json:"wrapperChecksumUrl"`
}
