package infrastructure

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/infrastructure/httpcache"
	"github.com/int128/gradleupdate/infrastructure/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"google.golang.org/appengine/urlfetch"
)

type GradleClient struct {
	dig.In
	ResponseCacheRepository gateways.ResponseCacheRepository
	Logger                  gateways.Logger
}

func (c *GradleClient) newClient(ctx context.Context) *http.Client {
	var transport http.RoundTripper
	transport = &urlfetch.Transport{Context: ctx}
	transport = &loggingTransport{Transport: transport, Name: "GradleClient", Logger: c.Logger}
	transport = &httpcache.Transport{Transport: transport, ResponseCacheRepository: c.ResponseCacheRepository, Logger: c.Logger}
	return &http.Client{Transport: transport}
}

func (c *GradleClient) GetCurrentVersion(ctx context.Context) (*infrastructure.CurrentVersionResponse, error) {
	client := c.newClient(ctx)
	req, err := http.NewRequest("GET", "https://services.gradle.org/versions/current", nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating a HTTP request")
	}
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting the current version from Gradle Service")
	}
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	var cvr infrastructure.CurrentVersionResponse
	if err := d.Decode(&cvr); err != nil {
		return nil, errors.Wrapf(err, "error while decoding JSON response from Gradle Service")
	}
	return &cvr, nil
}
