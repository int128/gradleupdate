package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/app/registry"
	"github.com/int128/gradleupdate/app/usecases"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type sendPullRequest struct {
	repositories registry.Repositories
}

func (h *sendPullRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]

	u := usecases.SendPullRequestForUpdate{
		Repository:  h.repositories.Repository(ctx),
		PullRequest: h.repositories.PullRequest(ctx),
		Branch:      h.repositories.Branch(ctx),
		Commit:      h.repositories.Commit(ctx),
	}
	if err := u.Do(ctx, owner, repo); err != nil {
		log.Errorf(ctx, "Error while usecases.SendPullRequestForUpdate: %+v", err)
		http.Error(w, err.Error(), 500)
	}
}
