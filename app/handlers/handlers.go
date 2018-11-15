package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

// New returns a handler for all paths.
func New() http.Handler {
	r := mux.NewRouter()
	r.Handle("/landing", &landing{}).Methods("POST")
	r.Handle("/{owner}/{repo}/status", &repository{}).Methods("GET")
	r.Handle("/{owner}/{repo}/status.svg", &badge{}).Methods("GET")
	r.Handle("/{owner}/{repo}/pull", &pullRequest{}).Methods("POST")
	return r
}

func baseURL(r *http.Request) string {
	scheme := "http"
	if r.Header.Get("X-AppEngine-Https") == "on" {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
