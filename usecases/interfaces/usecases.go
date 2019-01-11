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

//go:generate mockgen -destination mock_usecases/get_repo.go -package mock_usecases github.com/int128/gradleupdate/usecases/interfaces GetRepository

type GetRepository interface {
	Do(ctx context.Context, id domain.RepositoryID) (*GetRepositoryResponse, error)
}

type GetRepositoryResponse struct {
	Repository     domain.Repository
	CurrentVersion domain.GradleVersion
	LatestVersion  domain.GradleVersion
	UpToDate       bool
}

//go:generate mockgen -destination mock_usecases/send_pull_request.go -package mock_usecases github.com/int128/gradleupdate/usecases/interfaces SendPullRequest

type SendPullRequest interface {
	Do(ctx context.Context, id domain.RepositoryID) error
}
