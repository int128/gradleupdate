package httpcache

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
)

func TestConditionalRequestIfNoneMatch_CreateUpdate(t *testing.T) {
	s, h := setupTestServer(t)
	defer s.Close()
	cache := gateways.NewInMemoryCacheRepository()
	client := http.Client{
		Transport: &Transport{
			ResponseCacheRepository: cache,
			Transport:               http.DefaultTransport,
			Logger:                  gateways.NewLogger(t),
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
	cache := gateways.NewInMemoryCacheRepository()
	client := http.Client{
		Transport: &Transport{
			ResponseCacheRepository: cache,
			Transport:               http.DefaultTransport,
			Logger:                  gateways.NewLogger(t),
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
	cache := gateways.NewInMemoryCacheRepository()
	client := http.Client{
		Transport: &Transport{
			ResponseCacheRepository: cache,
			Transport:               http.DefaultTransport,
			Logger:                  gateways.NewLogger(t),
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
