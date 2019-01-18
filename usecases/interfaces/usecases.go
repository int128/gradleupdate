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

//go:generate mockgen -destination mock_usecases/request_update.go -package mock_usecases github.com/int128/gradleupdate/usecases/interfaces RequestUpdate

type RequestUpdate interface {
	Do(ctx context.Context, id domain.RepositoryID) error
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
