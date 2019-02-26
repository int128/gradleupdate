package usecases

import (
	"context"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
)

//go:generate mockgen -destination test_doubles/mock_usecases.go -package usecasesTestDoubles github.com/int128/gradleupdate/usecases/interfaces GetBadge,GetRepository,SendUpdate,BatchSendUpdates,SendPullRequest

type noSuchRepositoryErrorCauser interface {
	IsNoSuchRepositoryError(err error) bool
}

type GetBadge interface {
	IsNoGradleVersionError(err error) bool
	Do(ctx context.Context, id git.RepositoryID) (*GetBadgeResponse, error)
}

type GetBadgeResponse struct {
	CurrentVersion gradle.Version
	UpToDate       bool
}

type GetRepository interface {
	noSuchRepositoryErrorCauser
	Do(ctx context.Context, id git.RepositoryID) (*GetRepositoryResponse, error)
}

type GetRepositoryResponse struct {
	Repository                  git.Repository
	LatestGradleRelease         gradle.Release
	UpdatePreconditionViolation gradleupdate.PreconditionViolation
	UpdatePullRequestURL        git.PullRequestURL
}

type SendUpdate interface {
	noSuchRepositoryErrorCauser
	Do(ctx context.Context, id git.RepositoryID) error
	HasPreconditionViolation(err error) gradleupdate.PreconditionViolation
}

type BatchSendUpdates interface {
	Do(ctx context.Context) error
}

type SendPullRequest interface {
	Do(ctx context.Context, req SendPullRequestRequest) error
}

type SendPullRequestRequest struct {
	Base           git.RepositoryID
	HeadBranchName git.BranchName
	CommitMessage  string
	CommitFiles    []git.File
	Title          string
	Body           string
}
