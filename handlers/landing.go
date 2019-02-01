package handlers

import (
	"fmt"
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
		h.Logger.Infof(ctx, "Could not parse form: %s", err)
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}
	url := git.RepositoryURL(r.FormValue("url"))
	id := url.Parse()
	if id == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	to := fmt.Sprintf("/%s/%s/status", id.Owner, id.Name)
	http.Redirect(w, r, to, http.StatusFound)
}
