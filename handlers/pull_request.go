package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/usecases"
	"google.golang.org/appengine/log"
)

type SendPullRequest struct {
	ContextProvider          ContextProvider
	SendPullRequestForUpdate usecases.SendPullRequestForUpdate
}

func (h *SendPullRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := h.ContextProvider(r)
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]

	if err := h.SendPullRequestForUpdate.Do(ctx, owner, repo); err != nil {
		log.Errorf(ctx, "Error while usecases.SendPullRequestForUpdate: %+v", err)
		http.Error(w, err.Error(), 500)
	}
}
