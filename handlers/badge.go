package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/templates"
	"github.com/int128/gradleupdate/usecases"
	"google.golang.org/appengine/log"
)

type GetBadge struct {
	ContextProvider ContextProvider
	GetBadge        usecases.GetBadge
}

func (h *GetBadge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := h.ContextProvider(r)
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]

	badge, err := h.GetBadge.Do(ctx, owner, repo)
	switch {
	case err != nil:
		log.Warningf(ctx, "Could not get gradle wrapper version: %s", err)
		w.Header().Set("Content-Type", "image/svg+xml")
		templates.DarkBadge("unknown").WriteSVG(w)

	case badge.UpToDate:
		w.Header().Set("Content-Type", "image/svg+xml")
		templates.GreenBadge(string(badge.TargetVersion)).WriteSVG(w)

	case !badge.UpToDate:
		w.Header().Set("Content-Type", "image/svg+xml")
		templates.RedBadge(string(badge.TargetVersion)).WriteSVG(w)
	}
}
