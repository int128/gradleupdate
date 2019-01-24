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
}

func (c *GradleClient) newClient(ctx context.Context) *http.Client {
	var transport http.RoundTripper
	transport = &urlfetch.Transport{Context: ctx}
	transport = &loggingTransport{Transport: transport, Context: ctx, Name: "GradleClient"}
	transport = &httpcache.Transport{Transport: transport, Context: ctx, ResponseCacheRepository: c.ResponseCacheRepository}
	return &http.Client{Transport: transport}
}

func (c *GradleClient) GetCurrentVersion(ctx context.Context) (*infrastructure.CurrentVersionResponse, error) {
	client := c.newClient(ctx)
	resp, err := client.Get("https://services.gradle.org/versions/current")
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
