package handlers

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain/git"
	"go.uber.org/dig"
)

type RouterIn struct {
	dig.In
	Index            Index
	Landing          Landing
	GetRepository    GetRepository
	GetBadge         GetBadge
	SendUpdate       SendUpdate
	BatchSendUpdates BatchSendUpdates
}

func NewRouter(in RouterIn) *mux.Router {
	m := mux.NewRouter()
	m.Methods("GET").Path("/").Handler(&in.Index)
	m.Methods("POST").Path("/landing").Handler(&in.Landing)
	m.Methods("GET").Path("/{owner}/{repo}/status").Handler(&in.GetRepository)
	m.Methods("GET").Path("/{owner}/{repo}/status.svg").Handler(&in.GetBadge)
	m.Methods("POST").Path("/{owner}/{repo}/update").Handler(&in.SendUpdate)
	m.Methods("POST").Path("/internal/scan-updates").Handler(&in.BatchSendUpdates)
	return m
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
