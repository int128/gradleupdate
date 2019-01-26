package gateways

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/int128/gradleupdate/domain"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type GradleService struct {
	dig.In
	Client *http.Client
}

func (s *GradleService) GetCurrentVersion(ctx context.Context) (domain.GradleVersion, error) {
	req, err := http.NewRequest("GET", "https://services.gradle.org/versions/current", nil)
	if err != nil {
		return "", errors.Wrapf(err, "error while creating a HTTP request")
	}
	req = req.WithContext(ctx)
	resp, err := s.Client.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "error while getting the current version from Gradle Service")
	}
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	var cvr gradleCurrentVersionResponse
	if err := d.Decode(&cvr); err != nil {
		return "", errors.Wrapf(err, "error while decoding JSON response from Gradle Service")
	}
	return domain.GradleVersion(cvr.Version), nil
}

type gradleCurrentVersionResponse struct {
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
