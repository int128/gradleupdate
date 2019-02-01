package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/templates"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"go.uber.org/dig"
)

type SendUpdate struct {
	dig.In
	SendUpdate usecases.SendUpdate
	Logger     gateways.Logger
}

func (h *SendUpdate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]
	id := domain.RepositoryID{Owner: owner, Name: repo}

	if err := h.SendUpdate.Do(ctx, id); err != nil {
		h.Logger.Errorf(ctx, "error while sending a pull request for %s: %+v", id, err)
		w.Header().Set("content-type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		templates.WriteServerError(w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%s/status", owner, repo), http.StatusFound)
}
