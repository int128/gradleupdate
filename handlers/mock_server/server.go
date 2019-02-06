package main

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
)

func main() {
	c, err := newContainer()
	if err != nil {
		log.Fatalf("error while setting up a container: %+v", err)
	}
	if err := c.Invoke(run); err != nil {
		log.Fatalf("error while invoking app: %+v", err)
	}
}

func run(a app) error {
	m := http.NewServeMux()
	m.Handle("/", a.Router)

	static := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	m.Handle("/static/", static)

	a.Logger.Debugf(nil, "Open http://localhost:8080")
	if err := http.ListenAndServe("127.0.0.1:8080", m); err != nil {
		return errors.Wrapf(err, "error while listening on port")
	}
	return nil
}
