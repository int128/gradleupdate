package gateways

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

func NewInMemoryHTTPCacheRepository() *InMemoryHTTPCacheRepository {
	return &InMemoryHTTPCacheRepository{m: make(map[string][]byte)}
}

// InMemoryHTTPCacheRepository provides access to the response cache using a map.
type InMemoryHTTPCacheRepository struct {
	m map[string][]byte
	l sync.Mutex
}

func (r *InMemoryHTTPCacheRepository) ComputeKey(req *http.Request) string {
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

func (r *InMemoryHTTPCacheRepository) Find(ctx context.Context, k string, req *http.Request) (*http.Response, error) {
	r.l.Lock()
	v, ok := r.m[k]
	r.l.Unlock()
	if !ok {
		return nil, nil
	}
	b := bufio.NewReader(bytes.NewBuffer(v))
	resp, err := http.ReadResponse(b, req)
	if err != nil {
		return nil, errors.Wrapf(err, "could not decode response bytes")
	}
	return resp, nil
}

func (r *InMemoryHTTPCacheRepository) Save(ctx context.Context, k string, resp *http.Response) error {
	b, err := httputil.DumpResponse(resp, true) // DumpResponse preserves Body
	if err != nil {
		return errors.Wrapf(err, "could not dump response")
	}
	r.l.Lock()
	r.m[k] = b
	r.l.Unlock()
	return nil
}

func (r *InMemoryHTTPCacheRepository) Remove(ctx context.Context, k string) error {
	r.l.Lock()
	delete(r.m, k)
	r.l.Unlock()
	return nil
}

func (r *InMemoryHTTPCacheRepository) String() string {
	var b strings.Builder
	for k, v := range r.m {
		b.WriteString(fmt.Sprintf("%s=[%d]\n", k, len(v)))
	}
	return b.String()
}
