package usecases

import (
	"context"

	"github.com/int128/gradleupdate/domain"
)

//go:generate mockgen -destination mock_usecases/get_badge.go -package mock_usecases github.com/int128/gradleupdate/usecases/interfaces GetBadge

type GetBadge interface {
	Do(ctx context.Context, id domain.RepositoryID) (*GetBadgeResponse, error)
}

type GetBadgeResponse struct {
	CurrentVersion domain.GradleVersion
	UpToDate       bool
}

//go:generate mockgen -destination mock_usecases/get_repo.go -package mock_usecases github.com/int128/gradleupdate/usecases/interfaces GetRepository,GetRepositoryError

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

//go:generate mockgen -destination mock_usecases/send_update.go -package mock_usecases github.com/int128/gradleupdate/usecases/interfaces SendUpdate

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

//go:generate mockgen -destination mock_usecases/send_pull_request.go -package mock_usecases github.com/int128/gradleupdate/usecases/interfaces SendPullRequest

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
