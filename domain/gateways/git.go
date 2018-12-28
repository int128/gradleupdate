package gateways

import (
	"context"

	"github.com/int128/gradleupdate/domain"
)

type ForkBranchRequest struct {
	Base           domain.BranchIdentifier
	HeadBranchName string
	CommitMessage  string
	Files          []domain.File
}

type GitService interface {
	ForkBranch(ctx context.Context, req ForkBranchRequest) (*domain.Branch, error)
}
