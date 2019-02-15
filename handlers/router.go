package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/handlers/interfaces"
	"github.com/int128/gradleupdate/infrastructure"
	"go.uber.org/dig"
)

type RouterIn struct {
	dig.In
	Index                 Index
	Landing               Landing
	GetRepository         GetRepository
	GetBadge              GetBadge
	SendUpdate            SendUpdate
	TaskSendUpdate        TaskSendUpdate
	BatchSendUpdates      BatchSendUpdates
	CSRFMiddlewareFactory infrastructure.CSRFMiddlewareFactory
}

func NewRouter(in RouterIn) handlers.Router {
	r := mux.NewRouter()
	r.Methods("POST").Path("/internal/updates").Handler(&in.BatchSendUpdates)
	r.Methods("POST").Path("/internal/{owner}/{repo}/update").Handler(&in.TaskSendUpdate).Name("TaskSendUpdate")

	p := r.PathPrefix("/").Subrouter()
	p.Use(in.CSRFMiddlewareFactory.New())
	p.Methods("GET").Path("/").Handler(&in.Index)
	p.Methods("POST").Path("/landing").Handler(&in.Landing)
	p.Methods("GET").Path("/{owner}/{repo}/status").Handler(&in.GetRepository).Name("GetRepository")
	p.Methods("GET").Path("/{owner}/{repo}/status.svg").Handler(&in.GetBadge).Name("GetBadge")
	p.Methods("POST").Path("/{owner}/{repo}/update").Handler(&in.SendUpdate).Name("SendUpdate")

	r.NotFoundHandler = notFoundHandler("")
	r.MethodNotAllowedHandler = genericErrorHandler(http.StatusMethodNotAllowed)
	return r
}

func NewRouteResolver() handlers.RouteResolver {
	nullRouter, ok := NewRouter(RouterIn{}).(*mux.Router)
	if !ok {
		panic("NewRouter should return *mux.Router")
	}
	return &routeResolver{router: nullRouter}
}

type routeResolver struct {
	router *mux.Router
}

func (r *routeResolver) resolve(name string, pairs ...string) string {
	url, err := r.router.Get(name).URL(pairs...)
	if err != nil {
		panic(err) // should be fixed by tests
	}
	return url.String()
}

func (r *routeResolver) TaskSendUpdate(id git.RepositoryID) string {
	return r.resolve("TaskSendUpdate", "owner", id.Owner, "repo", id.Name)
}

func (r *routeResolver) GetRepositoryURL(id git.RepositoryID) string {
	return r.resolve("GetRepository", "owner", id.Owner, "repo", id.Name)
}

func (r *routeResolver) GetBadgeURL(id git.RepositoryID) string {
	return r.resolve("GetBadge", "owner", id.Owner, "repo", id.Name)
}

func (r *routeResolver) SendUpdateURL(id git.RepositoryID) string {
	return r.resolve("SendUpdate", "owner", id.Owner, "repo", id.Name)
}
