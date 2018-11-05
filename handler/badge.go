package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/service"
	"github.com/int128/gradleupdate/template"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// Badge handles a request for a badge.
type Badge struct{}

func (h *Badge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]

	status, err := service.GetGradleWrapperStatus(ctx, owner, repo)
	switch{
	case err != nil:
		log.Warningf(ctx, "Could not get gradle wrapper version: %s", err)
		t := template.Badge{
			LeftText:  "Gradle",
			LeftFill:  "#555",
			LeftWidth: 47,
			RightText: "unknown",
			RightFill: "#9f9f9f",
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		if err := t.Render(w); err != nil {
			log.Errorf(ctx, "Error while rendering SVG template: %s", err)
		}

	case status.UpToDate:
		t := template.Badge{
			LeftText:  "Gradle",
			LeftFill:  "#555",
			LeftWidth: 47,
			RightText: string(status.TargetVersion),
			RightFill: "#4c1",
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		if err := t.Render(w); err != nil {
			log.Errorf(ctx, "Error while rendering SVG template: %s", err)
		}

	case !status.UpToDate:
		t := template.Badge{
			LeftText:  "Gradle",
			LeftFill:  "#555",
			LeftWidth: 47,
			RightText: string(status.TargetVersion),
			RightFill: "#e05d44",
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		if err := t.Render(w); err != nil {
			log.Errorf(ctx, "Error while rendering SVG template: %s", err)
		}
	}
}
