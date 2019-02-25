package gateways

import (
	"context"
	"net/http"
	"time"

	"github.com/int128/gradleupdate/domain/config"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
)

//go:generate mockgen -destination test_doubles/mock_gateways.go -package gatewaysTestDoubles github.com/int128/gradleupdate/gateways/interfaces BadgeLastAccessRepository,GetRepositoryQuery,SendUpdateQuery,RepositoryRepository,PullRequestRepository,GradleReleaseRepository,Credentials,Toggles,Queue

type RepositoryError interface {
	error
	NoSuchEntity() bool
	AlreadyExists() bool
}

type BadgeLastAccessRepository interface {
	Save(ctx context.Context, a gradleupdate.BadgeLastAccess) error
	FindBySince(ctx context.Context, since time.Time) ([]gradleupdate.BadgeLastAccess, error)
}

type GetRepositoryQuery interface {
	Do(ctx context.Context, in GetRepositoryQueryIn) (*GetRepositoryQueryOut, error)
}

type GetRepositoryQueryIn struct {
	Repository     git.RepositoryID
	HeadBranchName git.BranchName
}

type GetRepositoryQueryOut struct {
	Repository              git.Repository
	PullRequestURL          git.PullRequestURL // a pull request associated with the head branch
	Readme                  git.FileContent
	GradleWrapperProperties git.FileContent
}

type SendUpdateQuery interface {
	Get(ctx context.Context, in SendUpdateQueryIn) (*SendUpdateQueryOut, error)
	ForkRepository(ctx context.Context, id git.RepositoryID) (*git.RepositoryID, error)
	CreateBranch(ctx context.Context, branch NewBranch) error
	UpdateBranch(ctx context.Context, branch NewBranch, force bool) error
}

type SendUpdateQueryIn struct {
	Repository     git.RepositoryID
	HeadBranchName git.BranchName
}

type SendUpdateQueryOut struct {
	BaseRepository          git.RepositoryID
	BaseBranch              git.BranchID
	BaseCommitSHA           git.CommitSHA
	BaseTreeSHA             git.TreeSHA
	HeadBranch              *git.BranchID // the head branch if an associated pull request exists
	HeadParentCommitSHA     git.CommitSHA // a parent of head branch (empty if the head has 2 or more parents)
	Readme                  git.FileContent
	GradleWrapperProperties git.FileContent
}

type NewBranch struct {
	Branch          git.BranchID
	ParentCommitSHA git.CommitSHA
	ParentTreeSHA   git.TreeSHA
	CommitMessage   string
	CommitFiles     []git.File
}

type RepositoryRepository interface {
	GetFileContent(context.Context, git.RepositoryID, string) (git.FileContent, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pull git.PullRequest) (*git.PullRequest, error)
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

type Queue interface {
	EnqueueSendUpdate(ctx context.Context, id git.RepositoryID) error
}

type Logger interface {
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
}
