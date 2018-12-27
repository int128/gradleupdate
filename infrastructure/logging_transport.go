package infrastructure

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type loggingTransport struct {
	Name      string
	Transport http.RoundTripper
}

func (t loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.Transport.RoundTrip(req)
	if res != nil {
		ctx := appengine.NewContext(req)
		log.Debugf(ctx, "[%s] %d %s %s", t.Name, res.StatusCode, req.Method, req.URL)
	}
	return res, err
}
