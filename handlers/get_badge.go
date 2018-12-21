package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/presenters/templates"
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
	id := domain.RepositoryIdentifier{Owner: owner, Name: repo}

	resp, err := h.GetBadge.Do(ctx, id)
	switch {
	case err != nil:
		log.Warningf(ctx, "could not get a badge for repository %s: %s", id, err)
		w.Header().Set("Content-Type", "image/svg+xml")
		templates.DarkBadge("unknown").WriteSVG(w)

	case resp.UpToDate:
		w.Header().Set("Content-Type", "image/svg+xml")
		templates.GreenBadge(string(resp.TargetVersion)).WriteSVG(w)

	case !resp.UpToDate:
		w.Header().Set("Content-Type", "image/svg+xml")
		templates.RedBadge(string(resp.TargetVersion)).WriteSVG(w)
	}
}
