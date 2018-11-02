package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// New returns a handler for all paths.
func New() http.Handler {
	r := mux.NewRouter()
	r.Handle("/{owner}/{repo}/status.svg", &Badge{}).Methods("GET")
	r.Handle("/", &Home{}).Methods("GET")
	return r
}
