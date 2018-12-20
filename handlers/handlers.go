package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/registry"
)

// New returns a handler for all paths.
func New() http.Handler {
	repositories := registry.NewRepositoriesRegistry()

	r := mux.NewRouter()
	r.Handle("/landing", &landing{}).Methods("POST")
	r.Handle("/{owner}/{repo}/status", &getStatus{repositories}).Methods("GET")
	r.Handle("/{owner}/{repo}/status.svg", &getBadge{repositories}).Methods("GET")
	r.Handle("/{owner}/{repo}/pull", &sendPullRequest{repositories}).Methods("POST")
	return r
}

func baseURL(r *http.Request) string {
	scheme := "http"
	if r.Header.Get("X-AppEngine-Https") == "on" {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
