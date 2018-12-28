package gateways

import (
	"context"

	"github.com/int128/gradleupdate/infrastructure"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
)

type PullRequestRepository struct {
	GitHubClient *infrastructure.GitHubClientFactory
}

func (r *PullRequestRepository) Query(ctx context.Context, q gateways.PullRequestQuery) ([]domain.PullRequest, error) {
	client := r.GitHubClient.New(ctx)
	payloads, _, err := client.PullRequests.List(ctx, q.Base.Repository.Owner, q.Base.Repository.Name, &github.PullRequestListOptions{
		Base:        q.Base.Name,
		Head:        q.Head.Name,
		State:       q.State,
		Direction:   q.Direction,
		Sort:        q.Sort,
		ListOptions: github.ListOptions{Page: q.Page, PerPage: q.PerPage},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "GitHub API returned error")
	}
	pulls := make([]domain.PullRequest, len(payloads))
	for i, payload := range payloads {
		head := payload.GetHead()
		base := payload.GetBase()
		pulls[i] = domain.PullRequest{
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
		}
	}
	return pulls, nil
}

func (r *PullRequestRepository) Create(ctx context.Context, pull domain.PullRequest) (*domain.PullRequest, error) {
	client := r.GitHubClient.New(ctx)
	payload, _, err := client.PullRequests.Create(ctx, pull.ID.Repository.Owner, pull.ID.Repository.Name, &github.NewPullRequest{
		Base:  github.String(pull.BaseBranch.Name),
		Head:  github.String(pull.HeadBranch.Name),
		Title: github.String(pull.Title),
		Body:  github.String(pull.Body),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "GitHub API returned error")
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

func (r *PullRequestRepository) Update(ctx context.Context, pull domain.PullRequest) (*domain.PullRequest, error) {
	client := r.GitHubClient.New(ctx)
	payload, _, err := client.PullRequests.Edit(ctx, pull.ID.Repository.Owner, pull.ID.Repository.Name, pull.ID.Number, &github.PullRequest{
		Title: github.String(pull.Title),
		Body:  github.String(pull.Body),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "GitHub API returned error")
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
