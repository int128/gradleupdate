// Package httpcache provides preliminary HTTP response cache by conditional requests.
package httpcache

import (
	"net/http"

	"github.com/int128/gradleupdate/domain/gateways"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type Transport struct {
	Transport               http.RoundTripper
	ResponseCacheRepository gateways.ResponseCacheRepository
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Transport == nil {
		return nil, errors.Errorf("given Transport is nil")
	}
	if t.ResponseCacheRepository == nil {
		return nil, errors.Errorf("given ResponseCacheRepository is nil")
	}
	if !isRequestCacheable(req) {
		return t.Transport.RoundTrip(req)
	}

	ctx := appengine.NewContext(req)
	cachedResp, err := t.ResponseCacheRepository.Find(ctx, req)
	if err != nil {
		log.Debugf(ctx, "error while finding response cache: %s", err)
	}
	if cachedResp == nil {
		resp, err := t.Transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		if isResponseCacheable(resp) {
			if err := t.ResponseCacheRepository.Save(ctx, req, resp); err != nil {
				log.Debugf(ctx, "error while saving response cache: %s", err)
			}
		}
		return resp, err
	}

	reqWithValidation := addCacheValidationHeaders(req, cachedResp)
	resp, err := t.Transport.RoundTrip(reqWithValidation)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotModified {
		return cachedResp, nil
	}
	if isResponseCacheable(resp) {
		if err := t.ResponseCacheRepository.Save(ctx, req, resp); err != nil {
			log.Debugf(ctx, "error while saving response cache: %s", err)
		}
	} else {
		if err := t.ResponseCacheRepository.Remove(ctx, req); err != nil {
			log.Debugf(ctx, "error while removing response cache: %s", err)
		}
	}
	return resp, nil
}

func isRequestCacheable(req *http.Request) bool {
	return req.Method == http.MethodGet
}

func isResponseCacheable(resp *http.Response) bool {
	return resp.Header.Get("etag") != ""
}

func addCacheValidationHeaders(req *http.Request, resp *http.Response) *http.Request {
	if resp.Header.Get("etag") == "" {
		return req
	}
	cloneReq := new(http.Request)
	*cloneReq = *req
	cloneReq.Header = make(http.Header)
	for key, values := range req.Header {
		cloneReq.Header[key] = values
	}
	cloneReq.Header.Set("if-none-match", resp.Header.Get("etag"))
	return cloneReq
}
