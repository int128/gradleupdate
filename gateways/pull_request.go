package gateways

import (
	"context"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type PullRequestRepository struct {
	dig.In
	Client *github.Client
}

func (r *PullRequestRepository) Create(ctx context.Context, pull git.PullRequest) (*git.PullRequest, error) {
	payload, _, err := r.Client.PullRequests.Create(ctx, pull.ID.Repository.Owner, pull.ID.Repository.Name, &github.NewPullRequest{
		Base:  github.String(pull.BaseBranch.Name),
		Head:  github.String(pull.HeadBranch.Repository.Owner + ":" + pull.HeadBranch.Name),
		Title: github.String(pull.Title),
		Body:  github.String(pull.Body),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error from GitHub API")
	}
	head := payload.GetHead()
	base := payload.GetBase()
	return &git.PullRequest{
		ID: git.PullRequestID{
			Repository: git.RepositoryID{Owner: base.GetUser().GetLogin(), Name: base.GetRepo().GetName()},
			Number:     payload.GetNumber(),
		},
		HeadBranch: git.BranchID{
			Repository: git.RepositoryID{Owner: head.GetUser().GetLogin(), Name: head.GetRepo().GetName()},
			Name:       head.GetRef(),
		},
		BaseBranch: git.BranchID{
			Repository: git.RepositoryID{Owner: base.GetUser().GetLogin(), Name: base.GetRepo().GetName()},
			Name:       base.GetRef(),
		},
		Title: payload.GetTitle(),
		Body:  payload.GetBody(),
	}, nil
}

func (r *PullRequestRepository) FindByBranch(ctx context.Context, baseBranch, headBranch git.BranchID) (*git.PullRequest, error) {
	pulls, _, err := r.Client.PullRequests.List(ctx, baseBranch.Repository.Owner, baseBranch.Repository.Name, &github.PullRequestListOptions{
		Base:        baseBranch.Name,
		Head:        headBranch.Repository.Owner + ":" + headBranch.Name,
		State:       "all",
		ListOptions: github.ListOptions{Page: 1, PerPage: 1},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error from GitHub API")
	}
	if len(pulls) > 1 {
		return nil, errors.Wrapf(err, "expect single pull request but got %+v", pulls)
	}
	if len(pulls) == 0 {
		return nil, nil
	}
	payload := pulls[0]
	head := payload.GetHead()
	base := payload.GetBase()
	return &git.PullRequest{
		ID: git.PullRequestID{
			Repository: git.RepositoryID{Owner: base.GetUser().GetLogin(), Name: base.GetRepo().GetName()},
			Number:     payload.GetNumber(),
		},
		HeadBranch: git.BranchID{
			Repository: git.RepositoryID{Owner: head.GetUser().GetLogin(), Name: head.GetRepo().GetName()},
			Name:       head.GetRef(),
		},
		BaseBranch: git.BranchID{
			Repository: git.RepositoryID{Owner: base.GetUser().GetLogin(), Name: base.GetRepo().GetName()},
			Name:       base.GetRef(),
		},
		Title: payload.GetTitle(),
		Body:  payload.GetBody(),
	}, nil
}
