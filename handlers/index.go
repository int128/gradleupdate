package handlers

import (
	"net/http"
	"time"

	"github.com/int128/gradleupdate/templates"
)

type Index struct{}

func (h *Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	w.Header().Set("cache-control", "public")
	w.Header().Set("expires", time.Now().Add(15*time.Second).Format(http.TimeFormat))
	templates.WriteIndex(w)
}
