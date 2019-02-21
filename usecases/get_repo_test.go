package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
	"github.com/int128/gradleupdate/domain/testdata"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
)

func TestGetRepository_Do(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}
	readmeContent := git.FileContent("![Gradle Status](https://gradleupdate.appspot.com/owner/repo/status.svg)")

	for _, c := range []struct {
		name                  string
		content               git.FileContent
		preconditionViolation gradleupdate.PreconditionViolation
	}{
		{
			name:                  "up-to-date",
			content:               testdata.GradleWrapperProperties50,
			preconditionViolation: gradleupdate.AlreadyHasLatestGradle,
		},
		{
			name:                  "out-of-date",
			content:               testdata.GradleWrapperProperties4102,
			preconditionViolation: gradleupdate.ReadyToUpdate,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			getRepositoryQuery := gatewaysTestDoubles.NewMockGetRepositoryQuery(ctrl)
			getRepositoryQuery.EXPECT().
				Do(ctx, gateways.GetRepositoryQueryIn{
					Repository:     repositoryID,
					HeadBranchName: gradleupdate.BranchFor(repositoryID.Owner, "5.0"),
				}).
				Return(&gateways.GetRepositoryQueryOut{
					Repository:              git.Repository{ID: repositoryID},
					Readme:                  readmeContent,
					GradleWrapperProperties: c.content,
					PullRequestURL:          "https://example.com/pull/1",
				}, nil)

			gradleService := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
			gradleService.EXPECT().
				GetCurrent(ctx).
				Return(&gradle.Release{Version: "5.0"}, nil)

			u := GetRepository{
				GetRepositoryQuery:      getRepositoryQuery,
				GradleReleaseRepository: gradleService,
			}
			resp, err := u.Do(ctx, repositoryID)
			if err != nil {
				t.Fatalf("could not do usecase: %s", err)
			}
			if resp.UpdatePreconditionViolation != c.preconditionViolation {
				t.Errorf("UpdatePreconditionViolation wants %v but %v", c.preconditionViolation, resp.UpdatePreconditionViolation)
			}
		})
	}

	t.Run("NoSuchRepository", func(t *testing.T) {
		getRepositoryQuery := gatewaysTestDoubles.NewMockGetRepositoryQuery(ctrl)
		getRepositoryQuery.EXPECT().
			Do(ctx, gateways.GetRepositoryQueryIn{
				Repository:     repositoryID,
				HeadBranchName: gradleupdate.BranchFor(repositoryID.Owner, "5.0"),
			}).
			Return(nil, &gatewaysTestDoubles.NoSuchEntityError{})

		gradleService := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleService.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		u := GetRepository{
			GetRepositoryQuery:      getRepositoryQuery,
			GradleReleaseRepository: gradleService,
		}
		resp, err := u.Do(ctx, repositoryID)
		if resp != nil {
			t.Errorf("resp wants nil but %+v", resp)
		}
		if err == nil {
			t.Fatalf("err wants non-nil but nil")
		}
		cause, ok := errors.Cause(err).(usecases.GetRepositoryError)
		if !ok {
			t.Fatalf("cause wants GetRepositoryError but %T", cause)
		}
		if cause.NoSuchRepository() != true {
			t.Errorf("NoSuchRepository wants true")
		}
	})
}
