package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradleupdate"
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
	id := git.RepositoryID{Owner: vars["owner"], Name: vars["repo"]}

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

	publicBadgeURL := gradleupdate.NewBadgeURL(id)
	publicRepositoryURL := gradleupdate.NewRepositoryURL(id)

	t := templates.Repository{
		Repository:                  resp.Repository,
		LatestGradleRelease:         resp.LatestGradleRelease,
		UpdatePreconditionViolation: resp.UpdatePreconditionViolation,
		BadgeMarkdown:               fmt.Sprintf(`[![Gradle Status](%s)](%s)`, publicBadgeURL, publicRepositoryURL),
		BadgeHTML:                   fmt.Sprintf(`<a href="%s"><img alt="Gradle Status" src="%s" /></a>`, publicRepositoryURL, publicBadgeURL),
		BadgeURL:                    resolveGetBadgeURL(id),
		RequestUpdateURL:            resolveSendUpdateURL(id),
	}
	t.WritePage(w, string(csrf.TemplateField(r)))
}
