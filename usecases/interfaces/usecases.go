package usecases

import (
	"context"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
)

//go:generate mockgen -destination test_doubles/mock_usecases.go -package usecases github.com/int128/gradleupdate/usecases/interfaces GetBadge,GetRepository,GetRepositoryError,SendUpdate,SendUpdateError,BatchSendUpdates,SendPullRequest

type GetBadge interface {
	Do(ctx context.Context, id git.RepositoryID) (*GetBadgeResponse, error)
}

type GetBadgeResponse struct {
	CurrentVersion gradle.Version
	UpToDate       bool
}

type GetRepository interface {
	Do(ctx context.Context, id git.RepositoryID) (*GetRepositoryResponse, error)
}

type GetRepositoryResponse struct {
	Repository                  git.Repository
	UpdatePreconditionViolation gradleupdate.PreconditionViolation
}

type GetRepositoryError interface {
	error
	NoSuchRepository() bool
}

type SendUpdate interface {
	Do(ctx context.Context, id git.RepositoryID) error
}

type SendUpdateError interface {
	error
	PreconditionViolation() gradleupdate.PreconditionViolation
}

type BatchSendUpdates interface {
	Do(ctx context.Context) error
}

type SendPullRequest interface {
	Do(ctx context.Context, req SendPullRequestRequest) error
}

type SendPullRequestRequest struct {
	Base           git.RepositoryID
	HeadBranchName string
	CommitMessage  string
	CommitFiles    []git.File
	Title          string
	Body           string
}
