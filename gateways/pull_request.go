package gateways

import (
	"context"

	"github.com/int128/gradleupdate/infrastructure"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain"
	"github.com/pkg/errors"
)

type PullRequestRepository struct {
	GitHubClientFactory *infrastructure.GitHubClientFactory
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
