package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handlers struct {
	Index         Index
	Landing       Landing
	GetRepository GetRepository
	GetBadge      GetBadge
	RequestUpdate RequestUpdate
}

func (h *Handlers) NewRouter() http.Handler {
	m := mux.NewRouter()
	m.Methods("GET").Path("/").Handler(&h.Index)
	m.Methods("POST").Path("/landing").Handler(&h.Landing)
	m.Methods("GET").Path("/{owner}/{repo}/status").Handler(&h.GetRepository)
	m.Methods("GET").Path("/{owner}/{repo}/status.svg").Handler(&h.GetBadge)
	m.Methods("POST").Path("/{owner}/{repo}/update").Handler(&h.RequestUpdate)
	return m
}

func baseURL(r *http.Request) string {
	scheme := "http"
	if r.Header.Get("X-AppEngine-Https") == "on" {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
