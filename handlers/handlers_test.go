package handlers_test

import (
	"testing"

	"github.com/int128/gradleupdate/gateways/interfaces"
	gatewaysTestDoubles "github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/int128/gradleupdate/handlers"
	"github.com/int128/gradleupdate/handlers/mock_server/di"
)

func invokeRouter(t *testing.T, runner func(handlers.Router)) {
	t.Helper()
	c, err := di.New()
	if err != nil {
		t.Fatalf("could not set up a container: %s", err)
	}
	if err := c.Provide(func() gateways.Logger {
		return gatewaysTestDoubles.NewLogger(t)
	}); err != nil {
		t.Fatalf("error while providing dependency: %+v", err)
	}
	if err := c.Invoke(runner); err != nil {
		t.Fatalf("error while running test: %+v", err)
	}
}
