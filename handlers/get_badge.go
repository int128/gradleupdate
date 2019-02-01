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
	owner, repo := vars["owner"], vars["repo"]
	id := git.RepositoryID{Owner: owner, Name: repo}

	w.Header().Set("content-type", "image/svg+xml")
	w.Header().Set("cache-control", "public")
	w.Header().Set("expires", time.Now().Add(15*time.Second).Format(http.TimeFormat))

	resp, err := h.GetBadge.Do(ctx, id)
	switch {
	case err != nil:
		h.Logger.Warnf(ctx, "could not get a badge for repository %s: %s", id, err)
		templates.DarkBadge("unknown").WriteSVG(w)

	case resp.UpToDate:
		templates.GreenBadge(string(resp.CurrentVersion)).WriteSVG(w)

	case !resp.UpToDate:
		templates.RedBadge(string(resp.CurrentVersion)).WriteSVG(w)
	}
}
