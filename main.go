package main

import (
	"net/http"

	"github.com/int128/gradleupdate/app/handlers"
	"google.golang.org/appengine"
)

func main() {
	http.Handle("/", handlers.New())
	appengine.Main()
}
