package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/infrastructure"
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
	r.Methods("POST").Path("/internal/scan-updates").Handler(&in.BatchSendUpdates)

	p := r.PathPrefix("/").Subrouter()
	p.Use(in.CSRFMiddlewareFactory.New())
	p.Methods("GET").Path("/").Handler(&in.Index)
	p.Methods("POST").Path("/landing").Handler(&in.Landing)
	p.Methods("GET").Path("/{owner}/{repo}/status").Handler(&in.GetRepository)
	p.Methods("GET").Path("/{owner}/{repo}/status.svg").Handler(&in.GetBadge)
	p.Methods("POST").Path("/{owner}/{repo}/update").Handler(&in.SendUpdate)
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
