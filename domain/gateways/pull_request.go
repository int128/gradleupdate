package gateways

import (
	"context"

	"github.com/int128/gradleupdate/domain"
)

//go:generate mockgen -destination mock_gateways/pull_request.go github.com/int128/gradleupdate/domain/gateways PullRequestRepository

type PullRequestRepository interface {
	Query(context.Context, PullRequestQuery) ([]domain.PullRequest, error)
	Create(context.Context, domain.PullRequest) (*domain.PullRequest, error)
	Update(context.Context, domain.PullRequest) (*domain.PullRequest, error)
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
