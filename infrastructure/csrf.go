package infrastructure

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"go.uber.org/dig"
)

type CSRFMiddlewareFactory struct {
	dig.In
	Credentials gateways.Credentials
	Logger      gateways.Logger
}

func (factory *CSRFMiddlewareFactory) New() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			credentials, err := factory.Credentials.Get(ctx)
			if err != nil {
				http.Error(w, "server error", http.StatusInternalServerError)
				factory.Logger.Errorf(ctx, "could not get credentials for CSRF middleware: %+v", err)
				return
			}
			m := csrf.Protect(credentials.CSRFKey, csrf.Secure(isHTTPS(r)))
			m(next).ServeHTTP(w, r)
		})
	}
}

func isHTTPS(r *http.Request) bool {
	return r.Header.Get("X-AppEngine-Https") == "on"
}
