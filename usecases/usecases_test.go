package usecases_test

import (
	"testing"

	"github.com/int128/gradleupdate/usecases"
)

func TestTimeProvider_Now(t *testing.T) {
	type exampleUsecase struct {
		usecases.TimeProvider
	}
	u := exampleUsecase{}
	now := u.Now()
	if now.IsZero() {
		t.Fatalf("Now wants non-zero but zero")
	}
}
