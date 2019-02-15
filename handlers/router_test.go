package handlers

import (
	"testing"

	"github.com/int128/gradleupdate/domain/git"
)

func TestNewRouteResolver(t *testing.T) {
	r := NewRouteResolver()

	t.Run("InternalSendUpdateURL", func(t *testing.T) {
		r.InternalSendUpdateURL(git.RepositoryID{Owner: "foo", Name: "bar"})
	})
	t.Run("GetRepositoryURL", func(t *testing.T) {
		r.GetRepositoryURL(git.RepositoryID{Owner: "foo", Name: "bar"})
	})
	t.Run("GetBadgeURL", func(t *testing.T) {
		r.GetBadgeURL(git.RepositoryID{Owner: "foo", Name: "bar"})
	})
	t.Run("SendUpdateURL", func(t *testing.T) {
		r.SendUpdateURL(git.RepositoryID{Owner: "foo", Name: "bar"})
	})
}
