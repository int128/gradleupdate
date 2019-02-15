package main

import (
	"log"
	"net/http"

	"github.com/int128/gradleupdate/di"
	"github.com/int128/gradleupdate/handlers/interfaces"
	"google.golang.org/appengine"
)

func run(router handlers.Router) {
	http.Handle("/", router)
	appengine.Main()
}

func main() {
	c, err := di.New()
	if err != nil {
		log.Fatalf("could not initialize the dependencies: %+v", err)
	}
	if err := c.Invoke(run); err != nil {
		log.Fatalf("could not run the application: %+v", err)
	}
}
