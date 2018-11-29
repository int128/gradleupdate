// Package httpcache provides preliminary HTTP response cache by conditional requests.
package httpcache

import (
	"bufio"
	"bytes"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httputil"
)

type Transport struct {
	Transport http.RoundTripper // Default to http.DefaultTransport
	Cache     Cache
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Cache == nil {
		return nil, errors.Errorf("Transport.Cache is nil")
	}
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	if !isRequestCacheable(req) {
		return transport.RoundTrip(req)
	}

	repository := responseCacheRepository{t.Cache}
	k := computeCacheKey(req)
	cachedResp, _ := repository.getByKey(k, req) /* ignore error */
	if cachedResp == nil {
		resp, err := transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		if isResponseCacheable(resp) {
			if err := repository.save(k, resp); err != nil { /* ignore error */
			}
		}
		return resp, err
	}

	req = addCacheValidationHeaders(req, cachedResp)
	resp, err := transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotModified {
		return cachedResp, nil
	}
	if isResponseCacheable(resp) {
		if err := repository.save(k, resp); err != nil { /* ignore error */
		}
	} else {
		if err := repository.delete(k); err != nil { /* ignore error */
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

type responseCacheRepository struct {
	cache Cache
}

func (r *responseCacheRepository) getByKey(k CacheKey, req *http.Request) (*http.Response, error) {
	v, err := r.cache.Get(k)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	b := bufio.NewReader(bytes.NewBuffer(v))
	resp, err := http.ReadResponse(b, req)
	if err != nil {
		return nil, errors.Wrapf(err, "Error while decoding response cache")
	}
	return resp, nil
}

func (r *responseCacheRepository) save(k CacheKey, resp *http.Response) error {
	b, err := httputil.DumpResponse(resp, true) // DumpResponse preserves Body
	if err != nil {
		return errors.Wrapf(err, "Error while encoding response cache")
	}
	v := CacheValue(b)
	return r.cache.Set(k, v)
}

func (r *responseCacheRepository) delete(k CacheKey) error {
	if err := r.cache.Delete(k); err != nil {
		return errors.Wrapf(err, "Error while deleting response cache")
	}
	return nil
}
