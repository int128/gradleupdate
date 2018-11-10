package main

import (
	"net/http"

	"github.com/int128/gradleupdate/handlers"
	"google.golang.org/appengine"
)

func main() {
	http.Handle("/", handlers.New())
	appengine.Main()
}
