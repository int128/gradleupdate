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

func TestSendUpdate_Do(t *testing.T) {
	ctx := context.Background()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}
	forkID := git.RepositoryID{Owner: "gradleupdate", Name: "repo"}
	readmeContent := git.FileContent("![Gradle Status](https://gradleupdate.appspot.com/owner/repo/status.svg)")

	t.Run("CreateBranchAndPullRequest", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gradleReleaseRepository := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleReleaseRepository.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		sendUpdateQuery := gatewaysTestDoubles.NewMockSendUpdateQuery(ctrl)
		sendUpdateQuery.EXPECT().
			Get(ctx, gateways.SendUpdateQueryIn{
				Repository:     repositoryID,
				HeadBranchName: "gradle-5.0-owner",
			}).
			Return(&gateways.SendUpdateQueryOut{
				BaseRepository:          repositoryID,
				BaseBranch:              git.BranchID{Repository: repositoryID, Name: "master"},
				BaseCommitSHA:           "COMMIT_SHA",
				BaseTreeSHA:             "TREE_SHA",
				HeadBranch:              nil, // indicates no pull request exists
				HeadParentCommitSHA:     "",
				Readme:                  readmeContent,
				GradleWrapperProperties: testdata.GradleWrapperProperties4102,
			}, nil)
		sendUpdateQuery.EXPECT().
			ForkRepository(ctx, repositoryID).
			Return(&forkID, nil)
		sendUpdateQuery.EXPECT().
			CreateBranch(ctx, gateways.NewBranch{
				Branch:          git.BranchID{Repository: forkID, Name: "gradle-5.0-owner"},
				ParentCommitSHA: "COMMIT_SHA",
				ParentTreeSHA:   "TREE_SHA",
				CommitMessage:   "Gradle 5.0",
				CommitFiles: []git.File{{
					Path:    gradle.WrapperPropertiesPath,
					Content: testdata.GradleWrapperProperties50,
				}},
			})

		pullRequestRepository := gatewaysTestDoubles.NewMockPullRequestRepository(ctrl)
		pullRequestRepository.EXPECT().
			Create(ctx, git.PullRequest{
				ID:         git.PullRequestID{Repository: repositoryID},
				BaseBranch: git.BranchID{Repository: repositoryID, Name: "master"},
				HeadBranch: git.BranchID{Repository: forkID, Name: "gradle-5.0-owner"},
				Title:      "Gradle 5.0",
				Body: `Gradle 5.0 is available.

This is sent by @gradleupdate. See https://gradleupdate.appspot.com/owner/repo/status for more.`,
			}).
			Return(&git.PullRequest{}, nil)

		u := SendUpdate{
			GradleReleaseRepository: gradleReleaseRepository,
			SendUpdateQuery:         sendUpdateQuery,
			PullRequestRepository:   pullRequestRepository,
			Logger:                  gatewaysTestDoubles.NewLogger(t),
		}
		err := u.Do(ctx, repositoryID)
		if err != nil {
			t.Fatalf("error while Do: %+v", err)
		}
	})

	t.Run("PullRequestExistsAndBranchIsUpToDate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gradleReleaseRepository := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleReleaseRepository.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		sendUpdateQuery := gatewaysTestDoubles.NewMockSendUpdateQuery(ctrl)
		sendUpdateQuery.EXPECT().
			Get(ctx, gateways.SendUpdateQueryIn{
				Repository:     repositoryID,
				HeadBranchName: "gradle-5.0-owner",
			}).
			Return(&gateways.SendUpdateQueryOut{
				BaseRepository:          repositoryID,
				BaseBranch:              git.BranchID{Repository: repositoryID, Name: "master"},
				BaseCommitSHA:           "COMMIT_SHA",
				BaseTreeSHA:             "TREE_SHA",
				HeadBranch:              &git.BranchID{Repository: forkID, Name: "gradle-5.0-owner"},
				HeadParentCommitSHA:     "COMMIT_SHA", // same as BaseCommitSHA
				Readme:                  readmeContent,
				GradleWrapperProperties: testdata.GradleWrapperProperties4102,
			}, nil)

		pullRequestRepository := gatewaysTestDoubles.NewMockPullRequestRepository(ctrl)

		u := SendUpdate{
			GradleReleaseRepository: gradleReleaseRepository,
			SendUpdateQuery:         sendUpdateQuery,
			PullRequestRepository:   pullRequestRepository,
			Logger:                  gatewaysTestDoubles.NewLogger(t),
		}
		err := u.Do(ctx, repositoryID)
		if err != nil {
			t.Fatalf("error while Do: %+v", err)
		}
	})

	t.Run("PullRequestExistsButBranchIsOutOfDate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gradleReleaseRepository := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleReleaseRepository.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		sendUpdateQuery := gatewaysTestDoubles.NewMockSendUpdateQuery(ctrl)
		sendUpdateQuery.EXPECT().
			Get(ctx, gateways.SendUpdateQueryIn{
				Repository:     repositoryID,
				HeadBranchName: "gradle-5.0-owner",
			}).
			Return(&gateways.SendUpdateQueryOut{
				BaseRepository:          repositoryID,
				BaseBranch:              git.BranchID{Repository: repositoryID, Name: "master"},
				BaseCommitSHA:           "COMMIT_SHA",
				BaseTreeSHA:             "TREE_SHA",
				HeadBranch:              &git.BranchID{Repository: forkID, Name: "gradle-5.0-owner"},
				HeadParentCommitSHA:     "OLD_COMMIT_SHA", // different from BaseCommitSHA
				Readme:                  readmeContent,
				GradleWrapperProperties: testdata.GradleWrapperProperties4102,
			}, nil)
		sendUpdateQuery.EXPECT().
			UpdateBranch(ctx, gateways.NewBranch{
				Branch:          git.BranchID{Repository: forkID, Name: "gradle-5.0-owner"},
				ParentCommitSHA: "COMMIT_SHA",
				ParentTreeSHA:   "TREE_SHA",
				CommitMessage:   "Gradle 5.0",
				CommitFiles: []git.File{{
					Path:    gradle.WrapperPropertiesPath,
					Content: testdata.GradleWrapperProperties50,
				}},
			}, true)

		pullRequestRepository := gatewaysTestDoubles.NewMockPullRequestRepository(ctrl)

		u := SendUpdate{
			GradleReleaseRepository: gradleReleaseRepository,
			SendUpdateQuery:         sendUpdateQuery,
			PullRequestRepository:   pullRequestRepository,
			Logger:                  gatewaysTestDoubles.NewLogger(t),
		}
		err := u.Do(ctx, repositoryID)
		if err != nil {
			t.Fatalf("error while Do: %+v", err)
		}
	})

	t.Run("AlreadyHasLatestGradle", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gradleReleaseRepository := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleReleaseRepository.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "4.10.2"}, nil)

		sendUpdateQuery := gatewaysTestDoubles.NewMockSendUpdateQuery(ctrl)
		sendUpdateQuery.EXPECT().
			Get(ctx, gateways.SendUpdateQueryIn{
				Repository:     repositoryID,
				HeadBranchName: "gradle-4.10.2-owner",
			}).
			Return(&gateways.SendUpdateQueryOut{
				BaseRepository:          repositoryID,
				BaseBranch:              git.BranchID{Repository: repositoryID, Name: "master"},
				BaseCommitSHA:           "COMMIT_SHA",
				BaseTreeSHA:             "TREE_SHA",
				HeadBranch:              nil, // indicates no pull request exists
				HeadParentCommitSHA:     "",
				Readme:                  readmeContent,
				GradleWrapperProperties: testdata.GradleWrapperProperties4102,
			}, nil)

		pullRequestRepository := gatewaysTestDoubles.NewMockPullRequestRepository(ctrl)

		u := SendUpdate{
			GradleReleaseRepository: gradleReleaseRepository,
			SendUpdateQuery:         sendUpdateQuery,
			PullRequestRepository:   pullRequestRepository,
			Logger:                  gatewaysTestDoubles.NewLogger(t),
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecases.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != gradleupdate.AlreadyHasLatestGradle {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.AlreadyHasLatestGradle, preconditionViolation)
		}
	})

	t.Run("NoGradleWrapperProperties", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gradleReleaseRepository := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleReleaseRepository.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		sendUpdateQuery := gatewaysTestDoubles.NewMockSendUpdateQuery(ctrl)
		sendUpdateQuery.EXPECT().
			Get(ctx, gateways.SendUpdateQueryIn{
				Repository:     repositoryID,
				HeadBranchName: "gradle-5.0-owner",
			}).
			Return(&gateways.SendUpdateQueryOut{
				BaseRepository:          repositoryID,
				BaseBranch:              git.BranchID{Repository: repositoryID, Name: "master"},
				BaseCommitSHA:           "COMMIT_SHA",
				BaseTreeSHA:             "TREE_SHA",
				HeadBranch:              nil, // indicates no pull request exists
				HeadParentCommitSHA:     "",
				Readme:                  readmeContent,
				GradleWrapperProperties: nil, // indicates no file
			}, nil)

		pullRequestRepository := gatewaysTestDoubles.NewMockPullRequestRepository(ctrl)

		u := SendUpdate{
			GradleReleaseRepository: gradleReleaseRepository,
			SendUpdateQuery:         sendUpdateQuery,
			PullRequestRepository:   pullRequestRepository,
			Logger:                  gatewaysTestDoubles.NewLogger(t),
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecases.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != gradleupdate.NoGradleWrapperProperties {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.NoGradleWrapperProperties, preconditionViolation)
		}
	})

	t.Run("NoGradleVersion", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gradleReleaseRepository := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleReleaseRepository.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		sendUpdateQuery := gatewaysTestDoubles.NewMockSendUpdateQuery(ctrl)
		sendUpdateQuery.EXPECT().
			Get(ctx, gateways.SendUpdateQueryIn{
				Repository:     repositoryID,
				HeadBranchName: "gradle-5.0-owner",
			}).
			Return(&gateways.SendUpdateQueryOut{
				BaseRepository:          repositoryID,
				BaseBranch:              git.BranchID{Repository: repositoryID, Name: "master"},
				BaseCommitSHA:           "COMMIT_SHA",
				BaseTreeSHA:             "TREE_SHA",
				HeadBranch:              nil, // indicates no pull request exists
				HeadParentCommitSHA:     "",
				Readme:                  readmeContent,
				GradleWrapperProperties: git.FileContent("INVALID"),
			}, nil)

		pullRequestRepository := gatewaysTestDoubles.NewMockPullRequestRepository(ctrl)

		u := SendUpdate{
			GradleReleaseRepository: gradleReleaseRepository,
			SendUpdateQuery:         sendUpdateQuery,
			PullRequestRepository:   pullRequestRepository,
			Logger:                  gatewaysTestDoubles.NewLogger(t),
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecases.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != gradleupdate.NoGradleVersion {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.NoGradleVersion, preconditionViolation)
		}
	})

	t.Run("NoReadme", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gradleReleaseRepository := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleReleaseRepository.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil).MaxTimes(1)

		sendUpdateQuery := gatewaysTestDoubles.NewMockSendUpdateQuery(ctrl)
		sendUpdateQuery.EXPECT().
			Get(ctx, gateways.SendUpdateQueryIn{
				Repository:     repositoryID,
				HeadBranchName: "gradle-5.0-owner",
			}).
			Return(&gateways.SendUpdateQueryOut{
				BaseRepository:          repositoryID,
				BaseBranch:              git.BranchID{Repository: repositoryID, Name: "master"},
				BaseCommitSHA:           "COMMIT_SHA",
				BaseTreeSHA:             "TREE_SHA",
				HeadBranch:              nil, // indicates no pull request exists
				HeadParentCommitSHA:     "",
				Readme:                  nil, // indicates no file
				GradleWrapperProperties: testdata.GradleWrapperProperties4102,
			}, nil)

		pullRequestRepository := gatewaysTestDoubles.NewMockPullRequestRepository(ctrl)

		u := SendUpdate{
			GradleReleaseRepository: gradleReleaseRepository,
			SendUpdateQuery:         sendUpdateQuery,
			PullRequestRepository:   pullRequestRepository,
			Logger:                  gatewaysTestDoubles.NewLogger(t),
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecases.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != gradleupdate.NoReadme {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.NoReadme, preconditionViolation)
		}
	})

	t.Run("NoReadmeBadge", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		gradleReleaseRepository := gatewaysTestDoubles.NewMockGradleReleaseRepository(ctrl)
		gradleReleaseRepository.EXPECT().
			GetCurrent(ctx).
			Return(&gradle.Release{Version: "5.0"}, nil)

		sendUpdateQuery := gatewaysTestDoubles.NewMockSendUpdateQuery(ctrl)
		sendUpdateQuery.EXPECT().
			Get(ctx, gateways.SendUpdateQueryIn{
				Repository:     repositoryID,
				HeadBranchName: "gradle-5.0-owner",
			}).
			Return(&gateways.SendUpdateQueryOut{
				BaseRepository:          repositoryID,
				BaseBranch:              git.BranchID{Repository: repositoryID, Name: "master"},
				BaseCommitSHA:           "COMMIT_SHA",
				BaseTreeSHA:             "TREE_SHA",
				HeadBranch:              nil, // indicates no pull request exists
				HeadParentCommitSHA:     "",
				Readme:                  git.FileContent("INVALID"),
				GradleWrapperProperties: testdata.GradleWrapperProperties4102,
			}, nil)

		pullRequestRepository := gatewaysTestDoubles.NewMockPullRequestRepository(ctrl)

		u := SendUpdate{
			GradleReleaseRepository: gradleReleaseRepository,
			SendUpdateQuery:         sendUpdateQuery,
			PullRequestRepository:   pullRequestRepository,
			Logger:                  gatewaysTestDoubles.NewLogger(t),
		}
		err := u.Do(ctx, repositoryID)
		if err == nil {
			t.Fatalf("error wants non-nil but nil")
		}
		sendUpdateError, ok := errors.Cause(err).(usecases.SendUpdateError)
		if !ok {
			t.Fatalf("cause wants SendUpdateError but %+v", errors.Cause(err))
		}
		preconditionViolation := sendUpdateError.PreconditionViolation()
		if preconditionViolation != gradleupdate.NoReadmeBadge {
			t.Errorf("PreconditionViolation wants %v but %v", gradleupdate.NoReadmeBadge, preconditionViolation)
		}
	})
}
