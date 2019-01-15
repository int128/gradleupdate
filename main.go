package main

import (
	"net/http"
	"os"

	"github.com/int128/gradleupdate/gateways"
	"github.com/int128/gradleupdate/handlers"
	"github.com/int128/gradleupdate/infrastructure"
	"github.com/int128/gradleupdate/usecases"
	"google.golang.org/appengine"
)

func main() {
	gitHubClientFactory := &infrastructure.GitHubClientFactory{
		Token:                   os.Getenv("GITHUB_TOKEN"),
		ResponseCacheRepository: &gateways.ResponseCacheRepository{},
	}
	gradleService := &gateways.GradleService{
		GradleClient: &infrastructure.GradleClient{
			ResponseCacheRepository: &gateways.ResponseCacheRepository{},
		},
	}
	h := handlers.Handlers{
		GetRepository: handlers.GetRepository{
			GetRepository: &usecases.GetRepository{
				GradleService:        gradleService,
				RepositoryRepository: &gateways.RepositoryRepository{GitHubClientFactory: gitHubClientFactory},
			},
		},
		GetBadge: handlers.GetBadge{
			GetBadge: &usecases.GetBadge{
				GradleService:             gradleService,
				RepositoryRepository:      &gateways.RepositoryRepository{GitHubClientFactory: gitHubClientFactory},
				BadgeLastAccessRepository: &gateways.BadgeLastAccessRepository{},
			},
		},
		RequestUpdate: handlers.RequestUpdate{
			RequestUpdate: &usecases.RequestUpdate{
				GradleService:        gradleService,
				RepositoryRepository: &gateways.RepositoryRepository{GitHubClientFactory: gitHubClientFactory},
				SendPullRequest: &usecases.SendPullRequest{
					RepositoryRepository:  &gateways.RepositoryRepository{GitHubClientFactory: gitHubClientFactory},
					PullRequestRepository: &gateways.PullRequestRepository{GitHubClientFactory: gitHubClientFactory},
					GitService:            &gateways.GitService{GitHubClientFactory: gitHubClientFactory},
				},
			},
		},
	}
	http.Handle("/", h.NewRouter())
	appengine.Main()
}
