package main

import (
	"context"
	"github.com/int128/gradleupdate/infrastructure/repositories"
	"github.com/int128/gradleupdate/usecases"
	"net/http"

	"github.com/int128/gradleupdate/handlers"
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
				Repository: &repositories.Repository{},
			},
		},
		GetBadge: handlers.GetBadge{
			ContextProvider: contextProvider,
			GetBadge: usecases.GetBadge{
				Repository:      &repositories.Repository{},
				BadgeLastAccess: &repositories.BadgeLastAccess{},
			},
		},
		SendPullRequest: handlers.SendPullRequest{
			ContextProvider: contextProvider,
			SendPullRequest: usecases.SendPullRequest{
				Repository:  &repositories.Repository{},
				PullRequest: &repositories.PullRequest{},
				Branch:      &repositories.Branch{},
				Commit:      &repositories.Commit{},
				Tree:        &repositories.Tree{},
			},
		},
	}
	http.Handle("/", h.NewRouter())
	appengine.Main()
}
