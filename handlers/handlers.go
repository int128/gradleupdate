package handlers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type ContextProvider func(*http.Request) context.Context

type Handlers struct {
	Index           Index
	Landing         Landing
	GetRepository   GetRepository
	GetBadge        GetBadge
	SendPullRequest SendPullRequest
}

func (h *Handlers) NewRouter() http.Handler {
	m := mux.NewRouter()
	m.Handle("/", &h.Index).Methods("GET")
	m.Handle("/landing", &h.Landing).Methods("POST")
	m.Handle("/{owner}/{repo}/status", &h.GetRepository).Methods("GET")
	m.Handle("/{owner}/{repo}/status.svg", &h.GetBadge).Methods("GET")
	m.Handle("/{owner}/{repo}/send-pull-request", &h.SendPullRequest).Methods("POST")
	return m
}

func baseURL(r *http.Request) string {
	scheme := "http"
	if r.Header.Get("X-AppEngine-Https") == "on" {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
