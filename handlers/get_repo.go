package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/templates"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type GetRepository struct {
	dig.In
	GetRepository usecases.GetRepository
	Logger        gateways.Logger
}

func (h *GetRepository) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]
	id := domain.RepositoryID{Owner: owner, Name: repo}

	resp, err := h.GetRepository.Do(ctx, id)
	if err != nil {
		if err, ok := errors.Cause(err).(usecases.GetRepositoryError); ok {
			switch {
			case err.NoSuchRepository():
				w.Header().Set("content-type", "text/html")
				w.WriteHeader(http.StatusNotFound)
				templates.WriteNotFoundError(w, fmt.Sprintf("no such a repository %s", id))
				return
			}
		}
		h.Logger.Errorf(ctx, "could not get the repository %s: %s", id, err)
		w.Header().Set("content-type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		templates.WriteServerError(w)
		return
	}

	w.Header().Set("content-type", "text/html")
	w.Header().Set("cache-control", "public")
	w.Header().Set("expires", time.Now().Add(15*time.Second).Format(http.TimeFormat))

	baseURL := baseURL(r)
	badgeURL := fmt.Sprintf("/%s/%s/status.svg", owner, repo)
	badgeFullURL := baseURL + badgeURL
	repositoryFullURL := baseURL + fmt.Sprintf("/%s/%s/status", owner, repo)

	t := templates.Repository{
		Repository:                  resp.Repository,
		GradleUpdatePreconditionOut: resp.GradleUpdatePreconditionOut,
		BadgeMarkdown:               fmt.Sprintf(`[![Gradle Status](%s)](%s)`, badgeFullURL, repositoryFullURL),
		BadgeHTML:                   fmt.Sprintf(`<a href="%s"><img alt="Gradle Status" src="%s" /></a>`, repositoryFullURL, badgeFullURL),
		BadgeURL:                    badgeURL,
	}
	t.WritePage(w)
}
