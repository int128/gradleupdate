package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/templates"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"google.golang.org/appengine/log"
)

type GetRepository struct {
	GetRepository usecases.GetRepository
}

func (h *GetRepository) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]
	id := domain.RepositoryID{Owner: owner, Name: repo}

	resp, err := h.GetRepository.Do(ctx, id)
	if err != nil {
		log.Warningf(ctx, "could not get the repository %s: %s", id, err)
		http.Error(w, "could not get the repository", 500)
		return
	}

	t := templates.Repository{
		Repository:       resp.Repository,
		CurrentVersion:   resp.CurrentVersion,
		LatestVersion:    resp.LatestVersion,
		UpToDate:         resp.UpToDate,
		ThisURL:          fmt.Sprintf("/%s/%s/status", owner, repo),
		BadgeURL:         fmt.Sprintf("/%s/%s/status.svg", owner, repo),
		RequestUpdateURL: fmt.Sprintf("/%s/%s/update", owner, repo),
		BaseURL:          baseURL(r),
	}
	w.Header().Set("content-type", "text/html")
	w.Header().Set("cache-control", "public")
	w.Header().Set("expires", time.Now().Add(15*time.Second).Format(http.TimeFormat))
	t.WritePage(w)
}
