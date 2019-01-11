package main

import (
	"log"
	"net/http"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/templates"
)

const addr = "127.0.0.1:8080"

func main() {
	http.Handle("/static", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/Badge", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "image/svg+xml")
		templates.RedBadge("5.0").WriteSVG(w)
	})

	http.HandleFunc("/Repository", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		templates.Repository{
			Repository: domain.Repository{
				ID:          domain.RepositoryID{Owner: "int128", Name: "gradleupdate"},
				Description: "Automatic Gradle Update Service",
				AvatarURL:   "https://avatars0.githubusercontent.com/u/321266",
			},
			LatestVersion: "5.1",
			UpToDate:      false,
			ThisURL:       "/Repository",
			BadgeURL:      "/Badge",
			BaseURL:       "https://gradleupdate.appspot.com",
		}.WritePage(w)
	})

	log.Printf("Open http://%s", addr)
	if err := http.ListenAndServe(addr, http.DefaultServeMux); err != nil {
		log.Fatal(err)
	}
}
