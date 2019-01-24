package infrastructure

import (
	"context"
	"net/http"

	"github.com/int128/gradleupdate/gateways/interfaces"
)

type loggingTransport struct {
	Name      string
	Context   context.Context
	Transport http.RoundTripper
	Logger    gateways.Logger
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.Transport.RoundTrip(req)
	if res != nil {
		t.Logger.Debugf(t.Context, "[%s] %d %s %s", t.Name, res.StatusCode, req.Method, req.URL)
	}
	return res, err
}
