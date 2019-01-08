package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/templates"
	"github.com/int128/gradleupdate/usecases"
	"google.golang.org/appengine/log"
)

type GetRepository struct {
	ContextProvider ContextProvider
	GetRepository   usecases.GetRepository
}

func (h *GetRepository) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := h.ContextProvider(r)
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]
	id := domain.RepositoryID{Owner: owner, Name: repo}

	resp, err := h.GetRepository.Do(ctx, id)
	if err != nil {
		log.Warningf(ctx, "could not get the repository %s: %s", id, err)
		http.Error(w, "could not get the repository", 500)
		return
	}

	thisURL := fmt.Sprintf("%s/%s/%s/status", baseURL(r), owner, repo)
	badgeURL := fmt.Sprintf("%s/%s/%s/status.svg", baseURL(r), owner, repo)

	w.Header().Set("content-type", "text/html")
	w.Header().Set("cache-control", "public")
	w.Header().Set("expires", time.Now().Add(1*time.Minute).Format(http.TimeFormat))
	templates.WriteRepository(w,
		owner,
		repo,
		resp.Repository.Description,
		resp.Repository.AvatarURL,
		thisURL,
		badgeURL)
}
