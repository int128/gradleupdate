package handlers

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain"
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
	id := domain.RepositoryID{Owner: owner, Name: repo}

	w.Header().Set("content-type", "image/svg+xml")
	w.Header().Set("cache-control", "public")
	w.Header().Set("expires", time.Now().Add(1*time.Minute).Format(http.TimeFormat))

	resp, err := h.GetBadge.Do(ctx, id)
	switch {
	case err != nil:
		log.Warningf(ctx, "could not get a badge for repository %s: %s", id, err)
		templates.DarkBadge("unknown").WriteSVG(w)

	case resp.UpToDate:
		templates.GreenBadge(string(resp.CurrentVersion)).WriteSVG(w)

	case !resp.UpToDate:
		templates.RedBadge(string(resp.CurrentVersion)).WriteSVG(w)
	}
}
