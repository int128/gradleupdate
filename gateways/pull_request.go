package gateways

import (
	"context"

	"github.com/google/go-github/v24/github"
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
		Base:  github.String(pull.BaseBranch.Name.String()),
		Head:  github.String(pull.HeadBranch.Repository.Owner + ":" + pull.HeadBranch.Name.String()),
		Title: github.String(pull.Title),
		Body:  github.String(pull.Body),
	})
	if err != nil {
		if err, ok := err.(*github.ErrorResponse); ok {
			if err.Response.StatusCode == 422 {
				// GitHub does not return dedicated code on already existing, so catch all for now.
				return nil, errors.Wrapf(&repositoryError{error: err, alreadyExists: true}, "pull request already exists")
			}
		}
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
			Name:       git.BranchName(head.GetRef()),
		},
		BaseBranch: git.BranchID{
			Repository: git.RepositoryID{Owner: base.GetUser().GetLogin(), Name: base.GetRepo().GetName()},
			Name:       git.BranchName(base.GetRef()),
		},
		Title: payload.GetTitle(),
		Body:  payload.GetBody(),
	}, nil
}
