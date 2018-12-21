package gateways

import (
	"context"
	"encoding/json"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/infrastructure"
	"github.com/pkg/errors"
)

type currentVersionResponse struct {
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

type GradleService struct{}

func (s *GradleService) GetCurrentVersion(ctx context.Context) (domain.GradleVersion, error) {
	client := infrastructure.NewClient(ctx)
	resp, err := client.Get("https://services.gradle.org/versions/current")
	if err != nil {
		return "", errors.Wrapf(err, "error while getting the current version from Gradle Service")
	}
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	var cvr currentVersionResponse
	if err := d.Decode(&cvr); err != nil {
		return "", errors.Wrapf(err, "error while decoding JSON response from Gradle Service")
	}
	return domain.GradleVersion(cvr.Version), nil
}
