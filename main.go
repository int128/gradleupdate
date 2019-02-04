package main

import (
	"log"
	"net/http"

	"github.com/int128/gradleupdate/di"
	"google.golang.org/appengine"
)

func main() {
	if err := di.Invoke(func(app di.App) {
		http.Handle("/", app.Router)
		appengine.Main()
	}); err != nil {
		log.Fatalf("could not run the application: %+v", err)
	}
}
