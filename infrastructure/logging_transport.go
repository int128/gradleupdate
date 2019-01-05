package infrastructure

import (
	"context"
	"net/http"

	"google.golang.org/appengine/log"
)

type loggingTransport struct {
	Name      string
	Context   context.Context
	Transport http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.Transport.RoundTrip(req)
	if res != nil {
		log.Debugf(t.Context, "[%s] %d %s %s", t.Name, res.StatusCode, req.Method, req.URL)
	}
	return res, err
}
