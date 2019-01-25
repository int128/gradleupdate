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

func NewInMemoryCacheRepository() *InMemoryCacheRepository {
	return &InMemoryCacheRepository{m: make(map[string][]byte)}
}

// InMemoryCacheRepository provides access to the response cache using a map.
type InMemoryCacheRepository struct {
	m map[string][]byte
	l sync.Mutex
}

func (c *InMemoryCacheRepository) Find(ctx context.Context, req *http.Request) (*http.Response, error) {
	k := computeResponseCacheKey(req)
	c.l.Lock()
	v, ok := c.m[k]
	c.l.Unlock()
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

func (c *InMemoryCacheRepository) Save(ctx context.Context, req *http.Request, resp *http.Response) error {
	b, err := httputil.DumpResponse(resp, true) // DumpResponse preserves Body
	if err != nil {
		return errors.Wrapf(err, "could not dump response")
	}
	k := computeResponseCacheKey(req)
	c.l.Lock()
	c.m[k] = b
	c.l.Unlock()
	return nil
}

func (c *InMemoryCacheRepository) Remove(ctx context.Context, req *http.Request) error {
	k := computeResponseCacheKey(req)
	c.l.Lock()
	delete(c.m, k)
	c.l.Unlock()
	return nil
}

func (c *InMemoryCacheRepository) String() string {
	var b strings.Builder
	for k, v := range c.m {
		b.WriteString(fmt.Sprintf("%s=[%d]\n", k, len(v)))
	}
	return b.String()
}

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
