package handlers

import (
	"fmt"
	"net/http"

	"github.com/int128/gradleupdate/domain"
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
	url := domain.RepositoryURL(r.FormValue("url"))
	id := url.Parse()
	if id == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	to := fmt.Sprintf("/%s/%s/status", id.Owner, id.Name)
	http.Redirect(w, r, to, http.StatusFound)
}
