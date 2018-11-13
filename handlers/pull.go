package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/service"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type pullRequest struct{}

func (p *pullRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]

	if err := service.CreateOrUpdatePullRequestForGradleWrapper(ctx, owner, repo, "4.10.2"); err != nil {
		log.Errorf(ctx, "Error while CreateOrUpdatePullRequestForGradleWrapper: %+v", err)
		http.Error(w, err.Error(), 500)
	}
}
