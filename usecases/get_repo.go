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

type GetRepositoryIn struct {
	dig.In
	GetRepositoryQuery      gateways.GetRepositoryQuery
	GradleReleaseRepository gateways.GradleReleaseRepository
}

type GetRepository struct {
	GetRepositoryIn
	noSuchRepositoryErrorCauser
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
		if usecase.GetRepositoryQuery.IsNoSuchEntityError(err) {
			return nil, errors.WithStack(&noSuchRepositoryError{id})
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
