package repositories

import (
	"context"

	"github.com/int128/gradleupdate/app/domain"
)

type PullRequest interface {
	Query(context.Context, PullRequestQuery) ([]domain.PullRequest, error)
	Create(context.Context, domain.PullRequest) (domain.PullRequest, error)
	Update(context.Context, domain.PullRequest) (domain.PullRequest, error)
}

type PullRequestQuery struct {
	Head      domain.BranchIdentifier
	Base      domain.BranchIdentifier
	State     string
	Direction string
	Sort      string
	PerPage   int
	Page      int
}
