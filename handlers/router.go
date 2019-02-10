package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/infrastructure"
	"github.com/int128/gradleupdate/templates"
	"go.uber.org/dig"
)

type Router interface {
	http.Handler
}

type RouterIn struct {
	dig.In
	Index                 Index
	Landing               Landing
	GetRepository         GetRepository
	GetBadge              GetBadge
	SendUpdate            SendUpdate
	BatchSendUpdates      BatchSendUpdates
	CSRFMiddlewareFactory infrastructure.CSRFMiddlewareFactory
}

func NewRouter(in RouterIn) Router {
	r := mux.NewRouter()
	r.Methods("POST").Path("/internal/updates").Handler(&in.BatchSendUpdates)

	p := r.PathPrefix("/").Subrouter()
	p.Use(in.CSRFMiddlewareFactory.New())
	p.Methods("GET").Path("/").Handler(&in.Index)
	p.Methods("POST").Path("/landing").Handler(&in.Landing)
	p.Methods("GET").Path("/{owner}/{repo}/status").Handler(&in.GetRepository)
	p.Methods("GET").Path("/{owner}/{repo}/status.svg").Handler(&in.GetBadge)
	p.Methods("POST").Path("/{owner}/{repo}/update").Handler(&in.SendUpdate)

	r.NotFoundHandler = notFoundHandler("")
	r.MethodNotAllowedHandler = genericErrorHandler(http.StatusMethodNotAllowed)
	return r
}

func resolveGetRepositoryURL(id git.RepositoryID) string {
	return fmt.Sprintf("/%s/%s/status", id.Owner, id.Name)
}
func resolveGetBadgeURL(id git.RepositoryID) string {
	return fmt.Sprintf("/%s/%s/status.svg", id.Owner, id.Name)
}
func resolveSendUpdateURL(id git.RepositoryID) string {
	return fmt.Sprintf("/%s/%s/update", id.Owner, id.Name)
}

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
