package di_test

import (
	"testing"

	"github.com/int128/gradleupdate/di"
	"github.com/int128/gradleupdate/handlers"
)

func TestNew(t *testing.T) {
	c, err := di.New()
	if err != nil {
		t.Fatalf("could not initialize the dependencies: %+v", err)
	}
	if err := c.Invoke(func(handlers.Router) {}); err != nil {
		t.Fatalf("could not run the application: %+v", err)
	}
}
