package handler

import (
	"net/http"
	"regexp"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type landing struct{}

func (h *landing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if err := r.ParseForm(); err != nil {
		log.Infof(ctx, "Could not parse form: %s", err)
		http.Error(w, "Could not parse form", 400)
		return
	}
	url := r.FormValue("url")
	if url == "" {
		http.Redirect(w, r, "/", 302)
	}
	ownerSlashRepo := h.extractGitHubRepository(url)
	if ownerSlashRepo == "" {
		http.Redirect(w, r, "/", 302)
	}
	http.Redirect(w, r, ownerSlashRepo, 302)
}

var regexpGitHubURL = regexp.MustCompile(`/\w+/\w+$`)

func (h *landing) extractGitHubRepository(url string) string {
	return regexpGitHubURL.FindString(url)
}
