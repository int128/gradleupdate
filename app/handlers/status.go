package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/app/registry"
	"github.com/int128/gradleupdate/app/templates"
	"github.com/int128/gradleupdate/app/usecases"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type getStatus struct {
	repositories registry.Repositories
}

func (h *getStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]

	u := usecases.GetRepositoryAndStatus{
		Repository: h.repositories.Repository(ctx),
	}
	out, err := u.Do(ctx, owner, repo)
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
