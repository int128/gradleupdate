package handlers

import (
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type landing struct {
	routerHolder
}

func (h *landing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if err := r.ParseForm(); err != nil {
		log.Infof(ctx, "Could not parse form: %s", err)
		http.Error(w, "Could not parse form", 400)
		return
	}
	url := r.FormValue("url")
	owner, repo := h.extractOwnerAndRepo(url)
	to, err := h.router.Get("repository").URL("owner", owner, "repo", repo)
	if err != nil {
		log.Infof(ctx, "Could not determine URL for %s/%s: %s", owner, repo, err)
		http.Redirect(w, r, "/", 302)
		return
	}
	http.Redirect(w, r, to.String(), 302)
}

func (h *landing) extractOwnerAndRepo(url string) (string, string) {
	s := strings.Split(url, "/")
	if len(s) < 2 {
		return "", ""
	}
	return s[len(s)-2], s[len(s)-1]
}
