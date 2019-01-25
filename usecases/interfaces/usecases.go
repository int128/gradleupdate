package usecases

import (
	"context"

	"github.com/int128/gradleupdate/domain"
)

//go:generate mockgen -destination test_doubles/mock_usecases.go -package usecases github.com/int128/gradleupdate/usecases/interfaces GetBadge,GetRepository,GetRepositoryError,SendUpdate,SendUpdateError,BatchSendUpdates,SendPullRequest

type GetBadge interface {
	Do(ctx context.Context, id domain.RepositoryID) (*GetBadgeResponse, error)
}

type GetBadgeResponse struct {
	CurrentVersion domain.GradleVersion
	UpToDate       bool
}

type GetRepository interface {
	Do(ctx context.Context, id domain.RepositoryID) (*GetRepositoryResponse, error)
}

type GetRepositoryResponse struct {
	Repository     domain.Repository
	CurrentVersion domain.GradleVersion
	LatestVersion  domain.GradleVersion
	UpToDate       bool
}

type GetRepositoryError interface {
	error
	NoSuchRepository() bool
	NoGradleVersion() bool
}

type SendUpdate interface {
	Do(ctx context.Context, id domain.RepositoryID, badgeURL string) error
}

type SendUpdateError interface {
	error
	NoGradleVersion() bool
	NoReadmeBadge() bool
	AlreadyHasLatestGradle() bool
}

type BatchSendUpdates interface {
	Do(ctx context.Context) error
}

type SendPullRequest interface {
	Do(ctx context.Context, req SendPullRequestRequest) error
}

type SendPullRequestRequest struct {
	Base           domain.RepositoryID
	HeadBranchName string
	CommitMessage  string
	CommitFiles    []domain.File
	Title          string
	Body           string
}
