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

type SendUpdate struct {
	dig.In
	SendUpdate    usecases.SendUpdate
	RouteResolver handlers.RouteResolver
	Logger        gateways.Logger
}

func (h *SendUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := git.RepositoryID{Owner: vars["owner"], Name: vars["repo"]}

	if err := h.SendUpdate.Do(ctx, id); err != nil {
		h.Logger.Errorf(ctx, "error while sending a pull request for %s: %+v", id, err)
		genericErrorHandler(http.StatusInternalServerError).ServeHTTP(w, r)
		return
	}

	url := h.RouteResolver.GetRepositoryURL(id)
	http.Redirect(w, r, url, http.StatusFound)
}
