package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/presenters/templates"
	"github.com/int128/gradleupdate/usecases"
	"google.golang.org/appengine/log"
)

type GetRepository struct {
	ContextProvider     ContextProvider
	GetRepositoryStatus usecases.GetRepository
}

func (h *GetRepository) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := h.ContextProvider(r)
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]

	out, err := h.GetRepositoryStatus.Do(ctx, owner, repo)
	if err != nil {
		log.Warningf(ctx, "Could not get the repository: %s/%s: %s", owner, repo, err)
		http.Error(w, "Could not get the repository", 500)
		return
	}

	thisURL := fmt.Sprintf("%s/%s/%s/status", baseURL(r), owner, repo)
	badgeURL := fmt.Sprintf("%s/%s/%s/status.svg", baseURL(r), owner, repo)

	w.Header().Set("content-type", "text/html")
	templates.WriteRepository(w,
		owner,
		repo,
		out.Repository.Description,
		out.Repository.AvatarURL,
		thisURL,
		badgeURL)
}
