package handlers_test

import (
	"net/http/httptest"
	"testing"

	"github.com/int128/gradleupdate/handlers/interfaces"
)

func TestGetBadge_ServeHTTP(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		invokeRouter(t, func(router handlers.Router) {
			r := httptest.NewRequest("GET", "/int128/example/status.svg", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			resp := w.Result()
			if resp.StatusCode != 200 {
				t.Errorf("StatusCode wants 200 but %v", resp.StatusCode)
			}
			contentType := resp.Header.Get("content-type")
			if w := "image/svg+xml"; contentType != w {
				t.Errorf("content-type wants %s but %s", w, contentType)
			}
		})
	})
	t.Run("NotFound", func(t *testing.T) {
		invokeRouter(t, func(router handlers.Router) {
			r := httptest.NewRequest("GET", "/foo/example/status.svg", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			resp := w.Result()
			if resp.StatusCode != 404 {
				t.Errorf("StatusCode wants 200 but %v", resp.StatusCode)
			}
			contentType := resp.Header.Get("content-type")
			if w := "image/svg+xml"; contentType != w {
				t.Errorf("content-type wants %s but %s", w, contentType)
			}
		})
	})
}
