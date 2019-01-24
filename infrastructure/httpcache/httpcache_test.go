package httpcache

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"sync"
	"testing"

	"github.com/int128/gradleupdate/gateways/testing_logger"
	"github.com/pkg/errors"
)

func TestConditionalRequestIfNoneMatch_CreateUpdate(t *testing.T) {
	s, h := setupTestServer(t)
	defer s.Close()
	cache := newMemoryCacheRepository()
	client := http.Client{
		Transport: &Transport{
			ResponseCacheRepository: cache,
			Transport:               http.DefaultTransport,
			Logger:                  testing_logger.New(t),
		},
	}
	req, err := http.NewRequest("GET", s.URL+"/target", nil)
	if err != nil {
		t.Fatalf("Error while creating new request: %s", err)
	}

	t.Run("CacheNotFound", func(t *testing.T) {
		h.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !assertMethodURL(t, r, http.MethodGet, "/target") {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if v := r.Header.Get("if-none-match"); v != "" {
				t.Errorf("if-none-match wants empty but %s", v)
			}
			w.Header().Set("etag", "ETAG1")
			if _, err := w.Write([]byte("foo")); err != nil {
				t.Errorf("Error while writing body")
			}
		})
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Client returned error: %s", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("StatusCode wants %d but %d", http.StatusOK, resp.StatusCode)
		}
		if body := readResponseBody(t, resp); body != "foo" {
			t.Errorf("Body wants %s but %s", "foo", body)
		}
	})

	t.Run("CacheHit", func(t *testing.T) {
		h.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !assertMethodURL(t, r, http.MethodGet, "/target") {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if v := r.Header.Get("if-none-match"); v != "ETAG1" {
				t.Errorf("if-none-match wants %s but %s", "ETAG1", v)
			}
			w.WriteHeader(http.StatusNotModified)
		})
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Client returned error: %s", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("StatusCode wants %d but %d", http.StatusOK, resp.StatusCode)
		}
		if body := readResponseBody(t, resp); body != "foo" {
			t.Errorf("Body wants %s but %s", "foo", body)
		}
	})

	t.Run("InvalidateCache", func(t *testing.T) {
		h.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !assertMethodURL(t, r, http.MethodGet, "/target") {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if v := r.Header.Get("if-none-match"); v != "ETAG1" {
				t.Errorf("if-none-match wants %s but %s", "ETAG1", v)
			}
			w.Header().Set("etag", "ETAG2") // issue a new etag
			if _, err := w.Write([]byte("hello")); err != nil {
				t.Errorf("Error while writing body")
			}
		})
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Client returned error: %s", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("StatusCode wants %d but %d", http.StatusOK, resp.StatusCode)
		}
		if body := readResponseBody(t, resp); body != "hello" {
			t.Errorf("Body wants %s but %s", "hello", body)
		}
	})

	t.Run("CacheHitWithNewETag", func(t *testing.T) {
		h.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !assertMethodURL(t, r, http.MethodGet, "/target") {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if v := r.Header.Get("if-none-match"); v != "ETAG2" {
				t.Errorf("if-none-match wants %s but %s", "ETAG2", v)
			}
			w.WriteHeader(http.StatusNotModified)
		})
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Client returned error: %s", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("StatusCode wants %d but %d", http.StatusOK, resp.StatusCode)
		}
		if body := readResponseBody(t, resp); body != "hello" {
			t.Errorf("Body wants %s but %s", "hello", body)
		}
	})
}

func TestConditionalRequestIfNoneMatch_CreateDelete(t *testing.T) {
	s, h := setupTestServer(t)
	defer s.Close()
	cache := newMemoryCacheRepository()
	client := http.Client{
		Transport: &Transport{
			ResponseCacheRepository: cache,
			Transport:               http.DefaultTransport,
			Logger:                  testing_logger.New(t),
		},
	}
	req, err := http.NewRequest("GET", s.URL+"/target", nil)
	if err != nil {
		t.Fatalf("Error while creating new request: %s", err)
	}

	t.Run("CacheNotFound", func(t *testing.T) {
		h.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !assertMethodURL(t, r, http.MethodGet, "/target") {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if v := r.Header.Get("if-none-match"); v != "" {
				t.Errorf("if-none-match wants empty but %s", v)
			}
			w.Header().Set("etag", "ETAG1")
			if _, err := w.Write([]byte("foo")); err != nil {
				t.Errorf("Error while writing body")
			}
		})
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Client returned error: %s", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("StatusCode wants %d but %d", http.StatusOK, resp.StatusCode)
		}
		if body := readResponseBody(t, resp); body != "foo" {
			t.Errorf("Body wants %s but %s", "foo", body)
		}
		t.Logf("cache: %s", cache)
	})

	t.Run("CacheHit", func(t *testing.T) {
		h.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !assertMethodURL(t, r, http.MethodGet, "/target") {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if v := r.Header.Get("if-none-match"); v != "ETAG1" {
				t.Errorf("if-none-match wants %s but %s", "ETAG1", v)
			}
			w.WriteHeader(http.StatusNotModified)
		})
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Client returned error: %s", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("StatusCode wants %d but %d", http.StatusOK, resp.StatusCode)
		}
		if body := readResponseBody(t, resp); body != "foo" {
			t.Errorf("Body wants %s but %s", "foo", body)
		}
		t.Logf("cache: %s", cache)
	})

	t.Run("InvalidateCache", func(t *testing.T) {
		h.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !assertMethodURL(t, r, http.MethodGet, "/target") {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if v := r.Header.Get("if-none-match"); v != "ETAG1" {
				t.Errorf("if-none-match wants %s but %s", "ETAG1", v)
			}
			if _, err := w.Write([]byte("hello")); err != nil {
				t.Errorf("Error while writing body")
			}
		})
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Client returned error: %s", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("StatusCode wants %d but %d", http.StatusOK, resp.StatusCode)
		}
		if body := readResponseBody(t, resp); body != "hello" {
			t.Errorf("Body wants %s but %s", "hello", body)
		}
		t.Logf("cache: %s", cache)
	})

	t.Run("Transparent", func(t *testing.T) {
		h.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !assertMethodURL(t, r, http.MethodGet, "/target") {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if v := r.Header.Get("if-none-match"); v != "" {
				t.Errorf("if-none-match wants empty but %s", v)
			}
			if _, err := w.Write([]byte("world")); err != nil {
				t.Errorf("Error while writing body")
			}
		})
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Client returned error: %s", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("StatusCode wants %d but %d", http.StatusOK, resp.StatusCode)
		}
		if body := readResponseBody(t, resp); body != "world" {
			t.Errorf("Body wants %s but %s", "world", body)
		}
		t.Logf("cache: %s", cache)
	})
}

func TestNotCacheableRequest(t *testing.T) {
	s, h := setupTestServer(t)
	defer s.Close()
	cache := newMemoryCacheRepository()
	client := http.Client{
		Transport: &Transport{
			ResponseCacheRepository: cache,
			Transport:               http.DefaultTransport,
			Logger:                  testing_logger.New(t),
		},
	}
	req, err := http.NewRequest("POST", s.URL+"/target", nil)
	if err != nil {
		t.Fatalf("Error while creating new request: %s", err)
	}
	h.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !assertMethodURL(t, r, http.MethodPost, "/target") {
			w.WriteHeader(http.StatusNotFound)
		}
		if v := r.Header.Get("if-none-match"); v != "" {
			t.Errorf("if-none-match wants empty but %s", v)
		}
		if _, err := w.Write([]byte("foo")); err != nil {
			t.Errorf("Error while writing body")
		}
	})
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Client returned error: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode wants %d but %d", http.StatusOK, resp.StatusCode)
	}
	if body := readResponseBody(t, resp); body != "foo" {
		t.Errorf("Body wants %s but %s", "foo", body)
	}
}

func assertMethodURL(t *testing.T, r *http.Request, method string, url string) bool {
	ok := true
	if r.Method != method {
		t.Errorf("request.Method wants %s but %s", method, r.Method)
		ok = false
	}
	if r.URL.String() != url {
		t.Errorf("request.URL wants %s but %s", url, r.URL)
		ok = false
	}
	return ok
}

func readResponseBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}
	return string(b)
}

type handlerHolder struct {
	http.Handler
}

func setupTestServer(t *testing.T) (*httptest.Server, *handlerHolder) {
	t.Helper()
	holder := handlerHolder{http.DefaultServeMux}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		holder.ServeHTTP(w, r)
		t.Logf("[testServer] %s %s", r.Method, r.URL)
	})
	return httptest.NewServer(handler), &holder
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

type memoryCacheRepository struct {
	m map[string][]byte
	l sync.Mutex
}

func newMemoryCacheRepository() *memoryCacheRepository {
	return &memoryCacheRepository{m: make(map[string][]byte)}
}

func (c *memoryCacheRepository) Find(ctx context.Context, req *http.Request) (*http.Response, error) {
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

func (c *memoryCacheRepository) Save(ctx context.Context, req *http.Request, resp *http.Response) error {
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

func (c *memoryCacheRepository) Remove(ctx context.Context, req *http.Request) error {
	k := computeResponseCacheKey(req)
	c.l.Lock()
	delete(c.m, k)
	c.l.Unlock()
	return nil
}

func (c *memoryCacheRepository) String() string {
	var b strings.Builder
	for k, v := range c.m {
		b.WriteString(fmt.Sprintf("%s=[%d]\n", k, len(v)))
	}
	return b.String()
}
