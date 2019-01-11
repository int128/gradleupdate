package main

import (
	"context"
	"net/http"
	"os"

	"github.com/int128/gradleupdate/gateways"
	"github.com/int128/gradleupdate/handlers"
	"github.com/int128/gradleupdate/infrastructure"
	"github.com/int128/gradleupdate/usecases"
	"google.golang.org/appengine"
)

func contextProvider(req *http.Request) context.Context {
	return appengine.NewContext(req)
}

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
		Landing: handlers.Landing{
			ContextProvider: contextProvider,
		},
		GetRepository: handlers.GetRepository{
			ContextProvider: contextProvider,
			GetRepository: &usecases.GetRepository{
				GradleService:        gradleService,
				RepositoryRepository: &gateways.RepositoryRepository{GitHubClientFactory: gitHubClientFactory},
			},
		},
		GetBadge: handlers.GetBadge{
			ContextProvider: contextProvider,
			GetBadge: &usecases.GetBadge{
				GradleService:             gradleService,
				RepositoryRepository:      &gateways.RepositoryRepository{GitHubClientFactory: gitHubClientFactory},
				BadgeLastAccessRepository: &gateways.BadgeLastAccessRepository{},
			},
		},
		SendPullRequest: handlers.SendPullRequest{
			ContextProvider: contextProvider,
			SendPullRequest: &usecases.SendPullRequest{
				GradleService:         gradleService,
				RepositoryRepository:  &gateways.RepositoryRepository{GitHubClientFactory: gitHubClientFactory},
				PullRequestRepository: &gateways.PullRequestRepository{GitHubClientFactory: gitHubClientFactory},
				GitService:            &gateways.GitService{GitHubClientFactory: gitHubClientFactory},
			},
		},
	}
	http.Handle("/", h.NewRouter())
	appengine.Main()
}
