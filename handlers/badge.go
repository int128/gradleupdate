package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/int128/gradleupdate/service"
	"github.com/int128/gradleupdate/templates"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type badge struct{}

func (h *badge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	vars := mux.Vars(r)
	owner, repo := vars["owner"], vars["repo"]

	status, err := service.GetGradleWrapperStatus(ctx, owner, repo)
	switch{
	case err != nil:
		log.Warningf(ctx, "Could not get gradle wrapper version: %s", err)
		t := templates.Badge{
			LeftText:  "Gradle",
			LeftFill:  "#555",
			LeftWidth: 47,
			RightText: "unknown",
			RightFill: "#9f9f9f",
			RightWidth: templates.BadgeTextWidth("unknown"),
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		t.WriteSVG(w)

	case status.UpToDate:
		t := templates.Badge{
			LeftText:  "Gradle",
			LeftFill:  "#555",
			LeftWidth: 47,
			RightText: string(status.TargetVersion),
			RightFill: "#4c1",
			RightWidth: templates.BadgeTextWidth(string(status.TargetVersion)),
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		t.WriteSVG(w)

	case !status.UpToDate:
		t := templates.Badge{
			LeftText:  "Gradle",
			LeftFill:  "#555",
			LeftWidth: 47,
			RightText: string(status.TargetVersion),
			RightFill: "#e05d44",
			RightWidth: templates.BadgeTextWidth(string(status.TargetVersion)),
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		t.WriteSVG(w)
	}
}
