package handlers

import (
	"net/http"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"go.uber.org/dig"
)

type Landing struct {
	dig.In
	Logger gateways.Logger
}

func (h *Landing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		h.Logger.Debugf(ctx, "error while parsing form: %+v", err)
		genericErrorHandler(http.StatusBadRequest).ServeHTTP(w, r)
		return
	}
	url := git.RepositoryURL(r.FormValue("url"))
	id := url.Parse()
	if id == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	repositoryURL := resolveGetRepositoryURL(*id)
	http.Redirect(w, r, repositoryURL, http.StatusFound)
}
