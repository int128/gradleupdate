package di_test

import (
	"testing"

	"github.com/int128/gradleupdate/di"
)

func TestInvoke(t *testing.T) {
	if err := di.Invoke(func(root di.App) {
		t.Logf("%+v", root)
	}); err != nil {
		t.Fatalf("error while di.Invoke: %+v", err)
	}
}
