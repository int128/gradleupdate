package handlers

import (
	"net/http"

	"github.com/int128/gradleupdate/domain/git"
)

type Router interface {
	http.Handler
}

type RouteResolver interface {
	TaskSendUpdate(id git.RepositoryID) string
	GetRepositoryURL(id git.RepositoryID) string
	GetBadgeURL(id git.RepositoryID) string
	SendUpdateURL(id git.RepositoryID) string
}
