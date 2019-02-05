package main

import (
	"testing"

	"github.com/gorilla/mux"
)

func Test_newContainer(t *testing.T) {
	c, err := newContainer()
	if err != nil {
		t.Fatalf("could not set up a container: %s", err)
	}
	if err := c.Invoke(func(r *mux.Router) {}); err != nil {
		t.Fatalf("could not resolve dependencies: %s", err)
	}
}
