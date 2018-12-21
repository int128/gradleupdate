package main

import (
	"context"
	"net/http"

	"github.com/int128/gradleupdate/gateways"
	"github.com/int128/gradleupdate/handlers"
	"github.com/int128/gradleupdate/usecases"
	"google.golang.org/appengine"
)

func contextProvider(req *http.Request) context.Context {
	return appengine.NewContext(req)
}

func main() {
	h := handlers.Handlers{
		Landing: handlers.Landing{
			ContextProvider: contextProvider,
		},
		GetRepository: handlers.GetRepository{
			ContextProvider: contextProvider,
			GetRepositoryStatus: usecases.GetRepository{
				GradleService:        &gateways.GradleService{},
				RepositoryRepository: &gateways.RepositoryRepository{},
			},
		},
		GetBadge: handlers.GetBadge{
			ContextProvider: contextProvider,
			GetBadge: usecases.GetBadge{
				GradleService:             &gateways.GradleService{},
				RepositoryRepository:      &gateways.RepositoryRepository{},
				BadgeLastAccessRepository: &gateways.BadgeLastAccessRepository{},
			},
		},
		SendPullRequest: handlers.SendPullRequest{
			ContextProvider: contextProvider,
			SendPullRequest: usecases.SendPullRequest{
				RepositoryRepository:  &gateways.RepositoryRepository{},
				PullRequestRepository: &gateways.PullRequestRepository{},
				Branch:                &gateways.Branch{},
				Commit:                &gateways.Commit{},
				Tree:                  &gateways.Tree{},
			},
		},
	}
	http.Handle("/", h.NewRouter())
	appengine.Main()
}
