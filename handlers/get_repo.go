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
	"github.com/int128/gradleupdate/handlers/interfaces"
	"github.com/int128/gradleupdate/templates"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type GetRepository struct {
	dig.In
	GetRepository usecases.GetRepository
	RouteResolver handlers.RouteResolver
	Logger        gateways.Logger
}

func (h *GetRepository) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := git.RepositoryID{Owner: vars["owner"], Name: vars["repo"]}

	resp, err := h.GetRepository.Do(ctx, id)
	if err != nil {
		if err, ok := errors.Cause(err).(usecases.GetRepositoryError); ok {
			if err.NoSuchRepository() {
				notFoundHandler(fmt.Sprintf("no such a repository %s", id)).ServeHTTP(w, r)
				return
			}
		}
		h.Logger.Errorf(ctx, "error while getting the repository %s: %+v", id, err)
		genericErrorHandler(http.StatusInternalServerError).ServeHTTP(w, r)
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
		UpdatePullRequestURL:        resp.UpdatePullRequestURL,
		BadgeMarkdown:               fmt.Sprintf(`[![Gradle Status](%s)](%s)`, publicBadgeURL, publicRepositoryURL),
		BadgeHTML:                   fmt.Sprintf(`<a href="%s"><img alt="Gradle Status" src="%s" /></a>`, publicRepositoryURL, publicBadgeURL),
		BadgeURL:                    h.RouteResolver.GetBadgeURL(id),
		RequestUpdateURL:            h.RouteResolver.SendUpdateURL(id),
	}
	t.WritePage(w, string(csrf.TemplateField(r)))
}
