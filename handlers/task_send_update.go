package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/handlers/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"go.uber.org/dig"
)

type TaskSendUpdate struct {
	dig.In
	SendUpdate    usecases.SendUpdate
	RouteResolver handlers.RouteResolver
	Logger        gateways.Logger
}

func (h *TaskSendUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := git.RepositoryID{Owner: vars["owner"], Name: vars["repo"]}

	if err := h.SendUpdate.Do(ctx, id); err != nil {
		h.Logger.Errorf(ctx, "error while sending a pull request for %s: %+v", id, err)
		// do not retry the task
	}
	w.WriteHeader(http.StatusOK)
}
