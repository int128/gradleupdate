package usecases

import (
	"context"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradleupdate"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type GetRepository struct {
	dig.In
	GetRepositoryQuery      gateways.GetRepositoryQuery
	GradleReleaseRepository gateways.GradleReleaseRepository
}

func (usecase *GetRepository) Do(ctx context.Context, id git.RepositoryID) (*usecases.GetRepositoryResponse, error) {
	release, err := usecase.GradleReleaseRepository.GetCurrent(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting the latest Gradle release")
	}
	out, err := usecase.GetRepositoryQuery.Do(ctx, gateways.GetRepositoryQueryIn{
		Repository:     id,
		HeadBranchName: gradleupdate.BranchFor(id.Owner, release.Version),
	})
	if err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			if err.NoSuchEntity() {
				return nil, errors.Wrapf(&getRepositoryError{error: err, noSuchRepository: true}, "repository %s not found", id)
			}
		}
		return nil, errors.Wrapf(err, "error while getting the repository %s", id)
	}
	precondition := gradleupdate.Precondition{
		BadgeURL:                gradleupdate.NewBadgeURL(id),
		LatestGradleRelease:     release,
		Readme:                  out.Readme,
		GradleWrapperProperties: out.GradleWrapperProperties,
	}
	return &usecases.GetRepositoryResponse{
		Repository:                  out.Repository,
		LatestGradleRelease:         *release,
		UpdatePreconditionViolation: gradleupdate.CheckPrecondition(precondition),
		UpdatePullRequestURL:        out.PullRequestURL,
	}, nil
}

type getRepositoryError struct {
	error
	noSuchRepository bool
}

func (err *getRepositoryError) NoSuchRepository() bool { return err.noSuchRepository }

var _ usecases.GetRepositoryError = &getRepositoryError{}
