package infrastructure

import (
	"github.com/google/go-github/v18/github"
)

func GitHubClient() *github.Client {
	return github.NewClient(nil)
}
