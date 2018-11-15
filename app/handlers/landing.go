package handlers

import (
	"fmt"
	"net/http"
	"strings"

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
	owner, repo := h.extractOwnerAndRepo(url)
	to := fmt.Sprintf("/%s/%s/status", owner, repo)
	http.Redirect(w, r, to, 302)
}

func (h *landing) extractOwnerAndRepo(url string) (string, string) {
	s := strings.Split(url, "/")
	if len(s) < 2 {
		return "", ""
	}
	return s[len(s)-2], s[len(s)-1]
}
