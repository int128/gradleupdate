package gateways

import (
	"context"
	"net/http"

	"github.com/int128/gradleupdate/domain"
)

type BadgeLastAccessRepository interface {
	Get(context.Context, domain.RepositoryID) (*domain.BadgeLastAccess, error)
	Put(context.Context, domain.BadgeLastAccess) error
}

type ForkBranchRequest struct {
	Base           domain.BranchID
	HeadBranchName string
	CommitMessage  string
	Files          []domain.File
}

type GitService interface {
	ForkBranch(ctx context.Context, req ForkBranchRequest) (*domain.Branch, error)
}

//go:generate mockgen -destination mock_gateways/gradle.go -package mock_gateways github.com/int128/gradleupdate/gateways/interfaces GradleService

type GradleService interface {
	GetCurrentVersion(ctx context.Context) (domain.GradleVersion, error)
}

//go:generate mockgen -destination mock_gateways/pull_request.go -package mock_gateways github.com/int128/gradleupdate/gateways/interfaces PullRequestRepository

type PullRequestRepository interface {
	Query(context.Context, PullRequestQuery) ([]domain.PullRequest, error)
	Create(context.Context, domain.PullRequest) (*domain.PullRequest, error)
	Update(context.Context, domain.PullRequest) (*domain.PullRequest, error)
}

type PullRequestQuery struct {
	Head      domain.BranchID
	Base      domain.BranchID
	State     string
	Direction string
	Sort      string
	PerPage   int
	Page      int
}

type RepositoryRepository interface {
	Get(context.Context, domain.RepositoryID) (*domain.Repository, error)
	GetFileContent(context.Context, domain.RepositoryID, string) (domain.FileContent, error)
	Fork(context.Context, domain.RepositoryID) (*domain.Repository, error)
}

type ResponseCacheRepository interface {
	Find(ctx context.Context, req *http.Request) (*http.Response, error)
	Save(ctx context.Context, req *http.Request, resp *http.Response) error
	Remove(ctx context.Context, req *http.Request) error
}
