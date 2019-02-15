package handlers

import (
	"net/http"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/handlers/interfaces"
	"go.uber.org/dig"
)

type Landing struct {
	dig.In
	RouteResolver handlers.RouteResolver
	Logger        gateways.Logger
}

func (h *Landing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		h.Logger.Debugf(ctx, "error while parsing form: %+v", err)
		genericErrorHandler(http.StatusBadRequest).ServeHTTP(w, r)
		return
	}

	id := git.RepositoryURL(r.FormValue("url")).Parse()
	if id == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	url := h.RouteResolver.GetRepositoryURL(*id)
	http.Redirect(w, r, url, http.StatusFound)
}
