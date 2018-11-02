package main

import (
	"net/http"

	"github.com/int128/gradleupdate/handler"
	"google.golang.org/appengine"
)

func main() {
	http.Handle("/", handler.New())
	appengine.Main()
}
