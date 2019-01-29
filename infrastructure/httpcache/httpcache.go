// Package httpcache provides preliminary HTTP response cache by conditional requests.
package httpcache

import (
	"net/http"

	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
)

type Transport struct {
	Transport           http.RoundTripper
	HTTPCacheRepository gateways.HTTPCacheRepository
	Logger              gateways.Logger
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Transport == nil {
		return nil, errors.Errorf("given Transport is nil")
	}
	if t.HTTPCacheRepository == nil {
		return nil, errors.Errorf("given HTTPCacheRepository is nil")
	}
	if !isRequestCacheable(req) {
		return t.Transport.RoundTrip(req)
	}

	ctx := req.Context()
	key := t.HTTPCacheRepository.ComputeKey(req)
	cachedResp, err := t.HTTPCacheRepository.Find(ctx, key, req)
	if err != nil {
		t.Logger.Debugf(ctx, "error while finding response cache: %s", err)
	}
	if cachedResp == nil {
		resp, err := t.Transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		if isResponseCacheable(resp) {
			if err := t.HTTPCacheRepository.Save(ctx, key, resp); err != nil {
				t.Logger.Debugf(ctx, "error while saving response cache: %s", err)
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
		if err := t.HTTPCacheRepository.Save(ctx, key, resp); err != nil {
			t.Logger.Debugf(ctx, "error while saving response cache: %s", err)
		}
	} else {
		if err := t.HTTPCacheRepository.Remove(ctx, key); err != nil {
			t.Logger.Debugf(ctx, "error while removing response cache: %s", err)
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
