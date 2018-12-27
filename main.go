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
	gitHubClient := &infrastructure.GitHubClient{
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
			GetRepository: usecases.GetRepository{
				GradleService:        gradleService,
				RepositoryRepository: &gateways.RepositoryRepository{GitHubClient: gitHubClient},
			},
		},
		GetBadge: handlers.GetBadge{
			ContextProvider: contextProvider,
			GetBadge: usecases.GetBadge{
				GradleService:             gradleService,
				RepositoryRepository:      &gateways.RepositoryRepository{GitHubClient: gitHubClient},
				BadgeLastAccessRepository: &gateways.BadgeLastAccessRepository{},
			},
		},
		SendPullRequest: handlers.SendPullRequest{
			ContextProvider: contextProvider,
			SendPullRequest: usecases.SendPullRequest{
				RepositoryRepository:  &gateways.RepositoryRepository{GitHubClient: gitHubClient},
				PullRequestRepository: &gateways.PullRequestRepository{GitHubClient: gitHubClient},
				Branch:                &gateways.Branch{GitHubClient: gitHubClient},
				Commit:                &gateways.Commit{GitHubClient: gitHubClient},
				Tree:                  &gateways.Tree{GitHubClient: gitHubClient},
			},
		},
	}
	http.Handle("/", h.NewRouter())
	appengine.Main()
}
