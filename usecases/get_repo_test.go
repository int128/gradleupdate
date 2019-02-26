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
				GetRepositoryIn: GetRepositoryIn{
					GetRepositoryQuery:      getRepositoryQuery,
					GradleReleaseRepository: gradleService,
				},
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
		noSuchEntityError := errors.New("NoSuchRepository")

		getRepositoryQuery := gatewaysTestDoubles.NewMockGetRepositoryQuery(ctrl)
		getRepositoryQuery.EXPECT().
			Do(ctx, gateways.GetRepositoryQueryIn{
				Repository:     repositoryID,
				HeadBranchName: gradleupdate.BranchFor(repositoryID.Owner, "5.0"),
			}).
			Return(nil, noSuchEntityError)
		getRepositoryQuery.EXPECT().
			IsNoSuchEntityError(noSuchEntityError).
			Return(true)

		gradleService := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleService.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		u := GetRepository{
			GetRepositoryIn: GetRepositoryIn{
				GetRepositoryQuery:      getRepositoryQuery,
				GradleReleaseRepository: gradleService,
			},
		}
		resp, err := u.Do(ctx, repositoryID)
		if resp != nil {
			t.Errorf("resp wants nil but %+v", resp)
		}
		if err == nil {
			t.Fatalf("err wants non-nil but nil")
		}
		if !u.IsNoSuchRepositoryError(err) {
			t.Errorf("IsNoSuchRepositoryError wants true but false")
		}
	})
}
