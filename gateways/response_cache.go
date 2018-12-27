package gateways

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"net/http"
	"net/http/httputil"

	"github.com/pkg/errors"
	"google.golang.org/appengine/memcache"
)

func computeResponseCacheKey(req *http.Request) string {
	var b bytes.Buffer
	for key, values := range req.Header {
		b.Write([]byte(key))
		for _, value := range values {
			b.Write([]byte(value))
		}
	}
	b.Write([]byte(req.Method))
	b.Write([]byte(req.URL.String()))
	h := sha512.Sum512(b.Bytes())
	e := base64.StdEncoding.EncodeToString(h[:])
	return e
}

type ResponseCacheRepository struct{}

func (r *ResponseCacheRepository) Find(ctx context.Context, req *http.Request) (*http.Response, error) {
	key := computeResponseCacheKey(req)
	item, err := memcache.Get(ctx, key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "could not get key %s from memcache", key)
	}
	b := bufio.NewReader(bytes.NewBuffer(item.Value))
	resp, err := http.ReadResponse(b, req)
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode response bytes")
	}
	return resp, nil
}

func (r *ResponseCacheRepository) Remove(ctx context.Context, req *http.Request) error {
	key := computeResponseCacheKey(req)
	err := memcache.Delete(ctx, key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil
		}
		return errors.Wrapf(err, "could not remove key %s from memcache", key)
	}
	return err
}

func (r *ResponseCacheRepository) Save(ctx context.Context, req *http.Request, resp *http.Response) error {
	b, err := httputil.DumpResponse(resp, true) // DumpResponse preserves Body
	if err != nil {
		return errors.Wrapf(err, "could not dump response")
	}
	key := computeResponseCacheKey(req)
	if err := memcache.Set(ctx, &memcache.Item{Key: key, Value: b}); err != nil {
		return errors.Wrapf(err, "could not save key %s into memcache", key)
	}
	return err
}
