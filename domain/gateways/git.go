package gateways

import (
	"context"

	"github.com/int128/gradleupdate/domain"
)

type Branch interface {
	Get(context.Context, domain.BranchIdentifier) (domain.Branch, error)
	Create(context.Context, domain.Branch) (domain.Branch, error)
	UpdateForce(context.Context, domain.Branch) (domain.Branch, error)
}

type Commit interface {
	Get(context.Context, domain.CommitIdentifier) (domain.Commit, error)
	Create(context.Context, domain.Commit) (domain.Commit, error)
}

type Tree interface {
	Create(context.Context, domain.RepositoryIdentifier, domain.TreeIdentifier, []domain.File) (domain.TreeIdentifier, error)
}
