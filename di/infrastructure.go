package di

import (
	"net/http"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/infrastructure"
)

var infrastructureDependencies = []interface{}{
	func(factory infrastructure.GitHubClientFactory) *github.Client {
		return factory.New()
	},
	func(factory infrastructure.HTTPClientFactory) *http.Client {
		return factory.New()
	},
}
