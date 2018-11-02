package handler

import (
	"net/http"

	"github.com/int128/gradleupdate/template"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// Badge handles a request for a badge.
type Badge struct{}

func (h *Badge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	t := template.Badge{
		LeftText: "Gradle",
		LeftFill: "#555",
		LeftWidth: 47,
		RightText: "0.0",
		RightFill: "#4c1",
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	if err := t.Render(w); err != nil {
		log.Errorf(ctx, "Error while rendering SVG template: %s", err)
		http.Error(w, "Error while rendering SVG template", 500)
	}
}
