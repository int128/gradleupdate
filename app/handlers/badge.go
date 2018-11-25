package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/app/registry"
	"github.com/int128/gradleupdate/app/templates"
	"github.com/int128/gradleupdate/app/usecases"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type getBadge struct {
	repositories registry.Repositories
}

func (h *getBadge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]

	u := usecases.GetStatus{
		Repository: h.repositories.Repository(ctx),
	}
	status, err := u.Do(ctx, owner, repo)
	switch {
	case err != nil:
		log.Warningf(ctx, "Could not get gradle wrapper version: %s", err)
		w.Header().Set("Content-Type", "image/svg+xml")
		templates.DarkBadge("unknown").WriteSVG(w)

	case status.UpToDate:
		w.Header().Set("Content-Type", "image/svg+xml")
		templates.GreenBadge(string(status.TargetVersion)).WriteSVG(w)

	case !status.UpToDate:
		w.Header().Set("Content-Type", "image/svg+xml")
		templates.RedBadge(string(status.TargetVersion)).WriteSVG(w)
	}
}
