package handlers

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/templates"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"go.uber.org/dig"
)

type GetBadge struct {
	dig.In
	GetBadge usecases.GetBadge
	Logger   gateways.Logger
}

func (h *GetBadge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := git.RepositoryID{Owner: vars["owner"], Name: vars["repo"]}

	w.Header().Set("content-type", "image/svg+xml")
	w.Header().Set("cache-control", "public")
	w.Header().Set("expires", time.Now().Add(15*time.Second).Format(http.TimeFormat))

	resp, err := h.GetBadge.Do(ctx, id)
	if err != nil {
		if h.GetBadge.IsNoGradleVersionError(err) {
			w.WriteHeader(http.StatusNotFound)
			templates.DarkBadge("unknown").WriteSVG(w)
			return
		}
		h.Logger.Errorf(ctx, "error while getting a badge for the repository %s: %+v", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		templates.DarkBadge("unknown").WriteSVG(w)
		return
	}
	if resp.UpToDate {
		templates.GreenBadge(string(resp.CurrentVersion)).WriteSVG(w)
		return
	}
	templates.RedBadge(string(resp.CurrentVersion)).WriteSVG(w)
}
