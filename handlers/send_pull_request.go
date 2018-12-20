package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
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

	if err := h.SendPullRequest.Do(ctx, owner, repo); err != nil {
		log.Errorf(ctx, "Error while usecases.SendPullRequest: %+v", err)
		http.Error(w, err.Error(), 500)
	}
}
