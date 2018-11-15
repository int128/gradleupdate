package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/app/service"
	"github.com/int128/gradleupdate/app/templates"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type repository struct {
	routerHolder
}

func (h *repository) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]

	ghr, err := service.GetRepository(ctx, owner, repo)
	if err != nil {
		log.Warningf(ctx, "Could not get the repository: %s/%s: %s", owner, repo, err)
		http.Error(w, "Could not get the repository", 500)
		return
	}
	thisURL, err := h.router.Get("repository").URL("owner", owner, "repo", repo)
	if err != nil {
		log.Warningf(ctx, "Could not resolve repository URL: %s", err)
		http.Error(w, "Internal Error", 500)
		return
	}
	badgeURL, err := h.router.Get("badge").URL("owner", owner, "repo", repo)
	if err != nil {
		log.Warningf(ctx, "Could not resolve badge URL: %s", err)
		http.Error(w, "Internal Error", 500)
		return
	}

	w.Header().Set("content-type", "text/html")
	templates.WriteRepository(w,
		owner,
		repo,
		ghr.GetDescription(),
		ghr.GetOwner().GetAvatarURL(),
		baseURL(r)+thisURL.String(),
		baseURL(r)+badgeURL.String())
}
