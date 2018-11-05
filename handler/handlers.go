package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// New returns a handler for all paths.
func New() http.Handler {
	r := mux.NewRouter()
	r.Handle("/landing", &landing{}).Methods("POST")
	r.Handle("/{owner}/{repo}", &repository{}).Methods("GET")
	r.Handle("/{owner}/{repo}/status.svg", &badge{}).Methods("GET")
	return r
}
