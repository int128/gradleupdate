package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/usecases"
	"google.golang.org/appengine/log"
)

type SendPullRequest struct {
	ContextProvider ContextProvider
	SendPullRequest usecases.SendPullRequest
}

func (h *SendPullRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := h.ContextProvider(r)
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]
	id := domain.RepositoryIdentifier{Owner: owner, Name: repo}

	if err := h.SendPullRequest.Do(ctx, id); err != nil {
		log.Errorf(ctx, "could not send a pull request for %s: %+v", id, err)
		http.Error(w, err.Error(), 500)
	}
}
