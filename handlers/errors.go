package handlers

import (
	"net/http"

	"github.com/int128/gradleupdate/templates"
)

func notFoundHandler(message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		templates.WriteNotFoundError(w, message)
	})
}

func genericErrorHandler(code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		w.WriteHeader(code)
		templates.WriteError(w)
	})
}
