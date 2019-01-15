package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"google.golang.org/appengine/log"
)

type RequestUpdate struct {
	RequestUpdate usecases.RequestUpdate
}

func (h *RequestUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]
	id := domain.RepositoryID{Owner: owner, Name: repo}

	if err := h.RequestUpdate.Do(ctx, id); err != nil {
		log.Errorf(ctx, "could not send a pull request for %s: %+v", id, err)
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%s/status", owner, repo), http.StatusFound)
}
