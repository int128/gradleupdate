package main

import (
	"net/http"

	"google.golang.org/appengine"
)

func handler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", 404)
}

func main() {
	http.HandleFunc("/", handler)
	appengine.Main()
}
