package infrastructure

import (
	"context"
	"net/http"

	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/infrastructure/httpcache"
	"go.uber.org/dig"
	"google.golang.org/appengine/urlfetch"
)

type HTTPClientFactory struct {
	dig.In
	ResponseCacheRepository gateways.ResponseCacheRepository
	Logger                  gateways.Logger
}

func (factory *HTTPClientFactory) New(ctx context.Context) *http.Client {
	var transport http.RoundTripper
	transport = &urlfetch.Transport{Context: ctx}
	transport = &loggingTransport{Transport: transport, Logger: factory.Logger}
	transport = &httpcache.Transport{Transport: transport, ResponseCacheRepository: factory.ResponseCacheRepository, Logger: factory.Logger}
	return &http.Client{Transport: transport}
}
