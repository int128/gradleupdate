package infrastructure

import (
	"net/http"

	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/infrastructure/httpcache"
	"go.uber.org/dig"
	"google.golang.org/appengine/urlfetch"
)

type HTTPClientFactory struct {
	dig.In
	ResponseCacheRepository gateways.HTTPCacheRepository
	Logger                  gateways.Logger
}

func (factory *HTTPClientFactory) New() *http.Client {
	var transport http.RoundTripper
	transport = &aeTransport{}
	transport = &loggingTransport{Transport: transport, Logger: factory.Logger}
	transport = &httpcache.Transport{Transport: transport, HTTPCacheRepository: factory.ResponseCacheRepository, Logger: factory.Logger}
	return &http.Client{Transport: transport}
}

type aeTransport struct{}

func (t *aeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	transport := &urlfetch.Transport{Context: ctx}
	return transport.RoundTrip(req)
}

type loggingTransport struct {
	Transport http.RoundTripper
	Logger    gateways.Logger
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	res, err := t.Transport.RoundTrip(req)
	if res != nil {
		t.Logger.Debugf(ctx, "%d %s %s", res.StatusCode, req.Method, req.URL)
	}
	return res, err
}
