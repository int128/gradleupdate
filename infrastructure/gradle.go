package infrastructure

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/int128/gradleupdate/domain/gateways"
	"github.com/int128/gradleupdate/infrastructure/httpcache"
	"github.com/pkg/errors"
	"google.golang.org/appengine/urlfetch"
)

type GradleClient struct {
	ResponseCacheRepository gateways.ResponseCacheRepository
}

func (c *GradleClient) newClient(ctx context.Context) *http.Client {
	var transport http.RoundTripper
	transport = &urlfetch.Transport{Context: ctx}
	transport = &loggingTransport{Transport: transport, Name: "GradleClient"}
	transport = &httpcache.Transport{Transport: transport, ResponseCacheRepository: c.ResponseCacheRepository}
	return &http.Client{Transport: transport}
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

func (c *GradleClient) GetCurrentVersion(ctx context.Context) (*CurrentVersionResponse, error) {
	client := c.newClient(ctx)
	resp, err := client.Get("https://services.gradle.org/versions/current")
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting the current version from Gradle Service")
	}
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	var cvr CurrentVersionResponse
	if err := d.Decode(&cvr); err != nil {
		return nil, errors.Wrapf(err, "error while decoding JSON response from Gradle Service")
	}
	return &cvr, nil
}
