package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

// New returns a handler for all paths.
func New() http.Handler {
	r := mux.NewRouter()
	rh := routerHolder{r}
	r.Handle("/landing", &landing{rh}).Methods("POST").Name("landing")
	r.Handle("/{owner}/{repo}/status", &repository{rh}).Methods("GET").Name("repository")
	r.Handle("/{owner}/{repo}/status.svg", &badge{}).Methods("GET").Name("badge")
	return r
}

type routerHolder struct {
	router *mux.Router
}

func baseURL(r *http.Request) string {
	scheme := "http"
	if r.Header.Get("X-AppEngine-Https") == "on" {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
