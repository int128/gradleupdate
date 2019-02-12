package gateways

import (
	"context"
	"net/http"
	"time"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/config"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
)

//go:generate mockgen -destination test_doubles/mock_gateways.go -package gatewaysTestDoubles github.com/int128/gradleupdate/gateways/interfaces BadgeLastAccessRepository,RepositoryLastUpdateRepository,RepositoryRepository,PullRequestRepository,GitService,GradleReleaseRepository,Credentials,Toggles

type RepositoryError interface {
	error
	NoSuchEntity() bool
	AlreadyExists() bool
}

type BadgeLastAccessRepository interface {
	Save(ctx context.Context, a gradleupdate.BadgeLastAccess) error
	FindBySince(ctx context.Context, since time.Time) ([]gradleupdate.BadgeLastAccess, error)
}

type RepositoryLastUpdateRepository interface {
	Save(ctx context.Context, a domain.RepositoryLastUpdate) error
}

type RepositoryRepository interface {
	Get(context.Context, git.RepositoryID) (*git.Repository, error)
	GetFileContent(context.Context, git.RepositoryID, string) (git.FileContent, error)
	GetReadme(ctx context.Context, id git.RepositoryID) (git.FileContent, error)
	Fork(context.Context, git.RepositoryID) (*git.Repository, error)
	GetBranch(ctx context.Context, branch git.BranchID) (*git.Branch, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pull git.PullRequest) (*git.PullRequest, error)
}

type GitService interface {
	CreateBranch(ctx context.Context, req PushBranchRequest) (*git.Branch, error)
	UpdateForceBranch(ctx context.Context, req PushBranchRequest) (*git.Branch, error)
}

type PushBranchRequest struct {
	BaseBranch    git.Branch
	HeadBranch    git.BranchID
	CommitMessage string
	CommitFiles   []git.File
}

type GradleReleaseRepository interface {
	GetCurrent(ctx context.Context) (*gradle.Release, error)
}

type Time interface {
	Now() time.Time
}

type Credentials interface {
	Get(ctx context.Context) (*config.Credentials, error)
}

type Toggles interface {
	Get(ctx context.Context) (*config.Toggles, error)
}

type HTTPCacheRepository interface {
	ComputeKey(req *http.Request) string
	Find(ctx context.Context, key string, req *http.Request) (*http.Response, error)
	Save(ctx context.Context, key string, resp *http.Response) error
	Remove(ctx context.Context, key string) error
}

type Logger interface {
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
}
