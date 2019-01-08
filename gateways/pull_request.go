package gateways

import (
	"context"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/infrastructure/interfaces"
	"github.com/pkg/errors"
)

type PullRequestRepository struct {
	GitHubClientFactory infrastructure.GitHubClientFactory
}

func (r *PullRequestRepository) Create(ctx context.Context, pull domain.PullRequest) (*domain.PullRequest, error) {
	client := r.GitHubClientFactory.New(ctx)
	payload, _, err := client.PullRequests.Create(ctx, pull.ID.Repository.Owner, pull.ID.Repository.Name, &github.NewPullRequest{
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
	return &domain.PullRequest{
		ID: domain.PullRequestID{
			Repository: domain.RepositoryID{Owner: base.GetUser().GetLogin(), Name: base.GetRepo().GetName()},
			Number:     payload.GetNumber(),
		},
		HeadBranch: domain.BranchID{
			Repository: domain.RepositoryID{Owner: head.GetUser().GetLogin(), Name: head.GetRepo().GetName()},
			Name:       head.GetRef(),
		},
		BaseBranch: domain.BranchID{
			Repository: domain.RepositoryID{Owner: base.GetUser().GetLogin(), Name: base.GetRepo().GetName()},
			Name:       base.GetRef(),
		},
		Title: payload.GetTitle(),
		Body:  payload.GetBody(),
	}, nil
}

func (r *PullRequestRepository) FindByBranch(ctx context.Context, baseBranch, headBranch domain.BranchID) (*domain.PullRequest, error) {
	client := r.GitHubClientFactory.New(ctx)
	pulls, _, err := client.PullRequests.List(ctx, baseBranch.Repository.Owner, baseBranch.Repository.Name, &github.PullRequestListOptions{
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
	return &domain.PullRequest{
		ID: domain.PullRequestID{
			Repository: domain.RepositoryID{Owner: base.GetUser().GetLogin(), Name: base.GetRepo().GetName()},
			Number:     payload.GetNumber(),
		},
		HeadBranch: domain.BranchID{
			Repository: domain.RepositoryID{Owner: head.GetUser().GetLogin(), Name: head.GetRepo().GetName()},
			Name:       head.GetRef(),
		},
		BaseBranch: domain.BranchID{
			Repository: domain.RepositoryID{Owner: base.GetUser().GetLogin(), Name: base.GetRepo().GetName()},
			Name:       base.GetRef(),
		},
		Title: payload.GetTitle(),
		Body:  payload.GetBody(),
	}, nil
}
