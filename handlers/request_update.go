package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/templates"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
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
	badgeURL := fmt.Sprintf("/%s/%s/status.svg", owner, repo)

	if err := h.RequestUpdate.Do(ctx, id, badgeURL); err != nil {
		if err, ok := errors.Cause(err).(usecases.RequestUpdateError); ok {
			switch {
			case err.NoBadgeInReadme():
				w.Header().Set("content-type", "text/html")
				w.WriteHeader(http.StatusNotFound)
				//TODO: replace with a dedicated error page
				templates.WriteNotFoundError(w, fmt.Sprintf("no badge URL (%s) found in the repository %s", badgeURL, id))
				return
			}
		}
		log.Errorf(ctx, "could not send a pull request for %s: %+v", id, err)
		w.Header().Set("content-type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		templates.WriteServerError(w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%s/status", owner, repo), http.StatusFound)
}
