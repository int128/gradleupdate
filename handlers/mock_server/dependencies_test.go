package main

import (
	"testing"
)

func Test_newContainer(t *testing.T) {
	c, err := newContainer()
	if err != nil {
		t.Fatalf("could not set up a container: %s", err)
	}
	if err := c.Invoke(func(app) {}); err != nil {
		t.Fatalf("could not resolve dependencies: %s", err)
	}
}
