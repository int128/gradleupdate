package usecases

import (
	"context"
	"fmt"

	"github.com/int128/gradleupdate/app/domain"
	"github.com/int128/gradleupdate/app/domain/repositories"
	"github.com/pkg/errors"
)

type SendPullRequestForUpdate struct {
	Repository  repositories.Repository
	Branch      repositories.Branch
	Commit      repositories.Commit
	PullRequest repositories.PullRequest
}

func (interactor *SendPullRequestForUpdate) Do(ctx context.Context, owner, repo string) error {
	latestRepository := domain.RepositoryIdentifier{Owner: "int128", Repo: "latest-gradle-wrapper"}
	targetRepository := domain.RepositoryIdentifier{Owner: owner, Repo: repo}
	files, err := interactor.downloadGradleWrapperFiles(ctx, latestRepository)
	if err != nil {
		return errors.Wrapf(err, "Could not find files of the latest Gradle wrapper")
	}
	base, err := interactor.Repository.Get(ctx, targetRepository)
	if err != nil {
		return errors.Wrapf(err, "Could not get the repository %s/%s", owner, repo)
	}
	head, err := interactor.Repository.Fork(ctx, targetRepository)
	if err != nil {
		return errors.Wrapf(err, "Could not fork the repository %s/%s", owner, repo)
	}
	version := "x.y.z" //TODO
	commit := domain.Commit{
		CommitIdentifier: domain.CommitIdentifier{RepositoryIdentifier: head.RepositoryIdentifier},
		Message:          fmt.Sprintf("Gradle %s", version),
	}
	headBranch := domain.BranchIdentifier{
		RepositoryIdentifier: head.RepositoryIdentifier,
		Branch:               fmt.Sprintf("gradle-%s-%s", version, owner),
	}
	if _, err := interactor.commitAndPush(ctx, base.DefaultBranch, headBranch, commit, files); err != nil {
		return errors.Wrapf(err, "Could not commit and push a branch %s", headBranch)
	}
	pull := domain.PullRequest{
		Head: headBranch,
		Base: base.DefaultBranch,
		PullRequestIdentifier: domain.PullRequestIdentifier{
			RepositoryIdentifier: head.RepositoryIdentifier,
		},
		Title: fmt.Sprintf("Gradle %s", version),
		Body: fmt.Sprintf(`This will upgrade the Gradle wrapper to the latest version %s.

This pull request is sent by @gradleupdate and based on [int128/latest-gradle-wrapper](https://github.com/int128/latest-gradle-wrapper).
`, version),
	}
	if _, err := interactor.openPullRequest(ctx, pull); err != nil {
		return errors.Wrapf(err, "Could not open a pull request %s", pull)
	}
	return nil
}

var gradleWrapperFiles = []domain.File{
	domain.File{
		Path: "gradle/wrapper/gradle-wrapper.properties",
		Mode: "100644",
	},
	domain.File{
		Path: "gradle/wrapper/gradle-wrapper.jar",
		Mode: "100644",
	},
	domain.File{
		Path: "gradlew",
		Mode: "100755",
	},
	domain.File{
		Path: "gradlew.bat",
		Mode: "100644",
	},
}

func (interactor *SendPullRequestForUpdate) downloadGradleWrapperFiles(ctx context.Context, id domain.RepositoryIdentifier) ([]domain.File, error) {
	r := make([]domain.File, len(gradleWrapperFiles))
	for i, f := range gradleWrapperFiles {
		resp, err := interactor.Repository.GetFile(ctx, id, f.Path)
		if err != nil {
			return nil, errors.Wrapf(err, "Could not get file %s", f.Path)
		}
		r[i] = f
		r[i].Content = resp.Content
	}
	return r, nil
}

func (interactor *SendPullRequestForUpdate) commitAndPush(ctx context.Context, base, head domain.BranchIdentifier, commit domain.Commit, files []domain.File) (domain.Branch, error) {
	headBranch, err := interactor.Branch.Get(ctx, head)
	if domain.IsNotFoundError(err) {
		baseBranch, err := interactor.Branch.Get(ctx, base)
		if err != nil {
			return domain.Branch{}, errors.Wrapf(err, "Could not get the base branch %s", base)
		}
		baseCommit, err := interactor.Commit.Get(ctx, baseBranch.Commit)
		if err != nil {
			return domain.Branch{}, errors.Wrapf(err, "Could not get the base commit %s", baseBranch.Commit)
		}
		commit.Parents = []domain.CommitIdentifier{baseCommit.CommitIdentifier}
		newHeadCommit, err := interactor.Commit.Create(ctx, commit, files)
		if err != nil {
			return domain.Branch{}, err
		}
		newHeadBranch, err := interactor.Branch.Create(ctx, domain.Branch{
			BranchIdentifier: head,
			Commit:           newHeadCommit.CommitIdentifier,
		})
		if err != nil {
			return domain.Branch{}, errors.Wrapf(err, "Could not create a branch %s", head)
		}
		return newHeadBranch, nil
	}
	if err != nil {
		return domain.Branch{}, errors.Wrapf(err, "Could not get the head branch %s", head)
	}

	baseBranch, err := interactor.Branch.Get(ctx, base)
	if err != nil {
		return domain.Branch{}, errors.Wrapf(err, "Could not get the base branch %s", base)
	}
	headCommit, err := interactor.Commit.Get(ctx, headBranch.Commit)
	if err != nil {
		return domain.Branch{}, errors.Wrapf(err, "Could not get the commit %s of head branch %s", headBranch.Commit, headBranch)
	}
	parent := headCommit.GetSingleParent()
	if parent != nil && parent.SHA == baseBranch.Commit.SHA {
		return headBranch, nil
	}
	baseCommit, err := interactor.Commit.Get(ctx, baseBranch.Commit)
	if err != nil {
		return domain.Branch{}, errors.Wrapf(err, "Could not get the base commit %s", baseBranch.Commit)
	}
	commit.Parents = []domain.CommitIdentifier{baseCommit.CommitIdentifier}
	newHeadCommit, err := interactor.Commit.Create(ctx, commit, files)
	if err != nil {
		return domain.Branch{}, err
	}
	newHeadBranch, err := interactor.Branch.UpdateForce(ctx, domain.Branch{
		BranchIdentifier: head,
		Commit:           newHeadCommit.CommitIdentifier,
	})
	if err != nil {
		return domain.Branch{}, errors.Wrapf(err, "Could not update the branch %s", head)
	}
	return newHeadBranch, nil
}

func (interactor *SendPullRequestForUpdate) openPullRequest(ctx context.Context, pull domain.PullRequest) (domain.PullRequest, error) {
	pulls, err := interactor.PullRequest.Query(ctx, repositories.PullRequestQuery{
		Head:      pull.Head,
		Base:      pull.Base,
		State:     "open",
		Direction: "desc",
		Sort:      "updated",
		Page:      1,
		PerPage:   1,
	})
	if err != nil {
		return domain.PullRequest{}, errors.Wrapf(err, "Could not find the pull request %s", pull.PullRequestIdentifier)
	}
	if len(pulls) > 1 {
		return domain.PullRequest{}, errors.Errorf("Expect single but got %d pull requests", len(pulls))
	}
	if len(pulls) == 1 {
		existent := pulls[0]
		existent.Body = pull.Body
		existent.Title = pull.Title
		updated, err := interactor.PullRequest.Update(ctx, existent)
		if err != nil {
			return domain.PullRequest{}, errors.Wrapf(err, "Could not update the pull request %s", pull.PullRequestIdentifier)
		}
		return updated, err
	}

	created, err := interactor.PullRequest.Create(ctx, pull)
	if err != nil {
		return domain.PullRequest{}, errors.Wrapf(err, "Could not create a pull request on the repository %s", pull.RepositoryIdentifier)
	}
	return created, nil
}
