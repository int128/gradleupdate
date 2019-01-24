package infrastructure

import (
	"net/http"

	"github.com/int128/gradleupdate/gateways/interfaces"
)

type loggingTransport struct {
	Name      string
	Transport http.RoundTripper
	Logger    gateways.Logger
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	res, err := t.Transport.RoundTrip(req)
	if res != nil {
		t.Logger.Debugf(ctx, "[%s] %d %s %s", t.Name, res.StatusCode, req.Method, req.URL)
	}
	return res, err
}
