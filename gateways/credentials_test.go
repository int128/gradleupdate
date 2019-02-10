package gateways

import (
	"os"
	"testing"

	"github.com/favclip/testerator"
	"github.com/go-test/deep"
	"github.com/int128/gradleupdate/domain/config"
	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"google.golang.org/appengine/datastore"
)

func TestNewCredentials(t *testing.T) {
	_, ctx, err := testerator.SpinUp()
	if err != nil {
		t.Fatalf("could not spin up appengine context: %s", err)
	}
	defer testerator.SpinDown()
	const base64csrfKey = "MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY="

	t.Run("FromEnv", func(t *testing.T) {
		runWithEnv(t, "GITHUB_TOKEN", "0123456789abcdef", func(t *testing.T) {
			runWithEnv(t, "CSRF_KEY", base64csrfKey, func(t *testing.T) {
				credentials := NewCredentials(gatewaysTestDoubles.NewLogger(t))
				c, err := credentials.Get(ctx)
				if err != nil {
					t.Fatalf("error while Get: %+v", err)
				}
				want := &config.Credentials{
					GitHubToken: "0123456789abcdef",
					CSRFKey:     []byte("0123456789abcdef0123456789abcdef"),
				}
				if diff := deep.Equal(want, c); diff != nil {
					t.Error(diff)
				}
			})
		})
	})

	t.Run("FromDatastore", func(t *testing.T) {
		credentials := NewCredentials(gatewaysTestDoubles.NewLogger(t))
		k := credentialsKey(ctx, "DEFAULT")
		if _, err := datastore.Put(ctx, k, &credentialsEntity{
			GitHubToken: "0123456789abcdef",
			CSRFKey:     base64csrfKey, // base64
		}); err != nil {
			t.Fatalf("error while putting an entity: %s", err)
		}
		c, err := credentials.Get(ctx)
		if err != nil {
			t.Fatalf("error while Get: %+v", err)
		}
		want := &config.Credentials{
			GitHubToken: "0123456789abcdef",
			CSRFKey:     []byte("0123456789abcdef0123456789abcdef"),
		}
		if diff := deep.Equal(want, c); diff != nil {
			t.Error(diff)
		}
	})
}

func runWithEnv(t *testing.T, k, v string, f func(t *testing.T)) {
	t.Helper()
	unset := func() {
		if err := os.Unsetenv(k); err != nil {
			t.Fatalf("error while unsetting environment variable %s: %s", k, err)
		}
	}
	defer unset()
	if err := os.Setenv(k, v); err != nil {
		t.Fatalf("error while setting environment variable %s=%s: %s", k, v, err)
	}
	f(t)
}
