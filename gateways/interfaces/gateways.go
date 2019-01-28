package gateways

import (
	"context"
	"net/http"
	"time"

	"github.com/int128/gradleupdate/domain"
)

//go:generate mockgen -destination test_doubles/mock_gateways.go -package gateways github.com/int128/gradleupdate/gateways/interfaces RepositoryError,BadgeLastAccessRepository,RepositoryLastScanRepository,RepositoryRepository,PullRequestRepository,GitService,GradleService

type RepositoryError interface {
	error
	NoSuchEntity() bool
}

type BadgeLastAccessRepository interface {
	Save(ctx context.Context, a domain.BadgeLastAccess) error
	FindBySince(ctx context.Context, since time.Time) ([]domain.BadgeLastAccess, error)
}

type RepositoryLastScanRepository interface {
	Save(ctx context.Context, a domain.RepositoryLastScan) error
}

type RepositoryRepository interface {
	Get(context.Context, domain.RepositoryID) (*domain.Repository, error)
	GetFileContent(context.Context, domain.RepositoryID, string) (domain.FileContent, error)
	GetReadme(ctx context.Context, id domain.RepositoryID) (domain.FileContent, error)
	Fork(context.Context, domain.RepositoryID) (*domain.Repository, error)
	GetBranch(ctx context.Context, branch domain.BranchID) (*domain.Branch, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pull domain.PullRequest) (*domain.PullRequest, error)
	FindByBranch(ctx context.Context, baseBranch, headBranch domain.BranchID) (*domain.PullRequest, error)
}

type GitService interface {
	CreateBranch(ctx context.Context, req PushBranchRequest) (*domain.Branch, error)
	UpdateForceBranch(ctx context.Context, req PushBranchRequest) (*domain.Branch, error)
}

type PushBranchRequest struct {
	BaseBranch    domain.Branch
	HeadBranch    domain.BranchID
	CommitMessage string
	CommitFiles   []domain.File
}

type GradleService interface {
	GetCurrentVersion(ctx context.Context) (domain.GradleVersion, error)
}

type ConfigRepository interface {
	Get(ctx context.Context) (*domain.Config, error)
}

type ResponseCacheRepository interface {
	Find(ctx context.Context, req *http.Request) (*http.Response, error)
	Save(ctx context.Context, req *http.Request, resp *http.Response) error
	Remove(ctx context.Context, req *http.Request) error
}

type Logger interface {
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
}
