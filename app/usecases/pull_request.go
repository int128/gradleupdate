package usecases

import (
	"context"
	"fmt"
	"google.golang.org/appengine/log"

	"github.com/int128/gradleupdate/app/domain"
	"github.com/int128/gradleupdate/app/domain/repositories"
	"github.com/pkg/errors"
)

type SendPullRequestForUpdate struct {
	Repository  repositories.Repository
	Branch      repositories.Branch
	Commit      repositories.Commit
	Tree        repositories.Tree
	PullRequest repositories.PullRequest
}

func (interactor *SendPullRequestForUpdate) Do(ctx context.Context, owner, repo string) error {
	latestRepository := domain.RepositoryIdentifier{Owner: "int128", Repo: "latest-gradle-wrapper"}
	targetRepository := domain.RepositoryIdentifier{Owner: owner, Repo: repo}

	files, err := interactor.downloadGradleWrapperFiles(ctx, latestRepository)
	if err != nil {
		return errors.Wrapf(err, "Could not find files of the latest Gradle wrapper")
	}
	version := findGradleWrapperVersion(files)
	if version == "" {
		return errors.Errorf("Could not find Gradle wrapper version from files %v", files)
	}
	log.Debugf(ctx, "Found Gradle wrapper %s in the repository %v", version, latestRepository)

	base, err := interactor.Repository.Get(ctx, targetRepository)
	if err != nil {
		return errors.Wrapf(err, "Could not get the repository %v", targetRepository)
	}
	baseBranch := base.DefaultBranch

	head, err := interactor.Repository.Fork(ctx, targetRepository)
	if err != nil {
		return errors.Wrapf(err, "Could not fork the repository %v", targetRepository)
	}
	log.Debugf(ctx, "Forked the repository %v into %v", targetRepository, head.RepositoryIdentifier)

	commit := domain.Commit{
		CommitIdentifier: domain.CommitIdentifier{Repository: head.RepositoryIdentifier},
		Message:          fmt.Sprintf("Gradle %s", version),
	}
	headBranch := domain.BranchIdentifier{
		Repository: head.RepositoryIdentifier,
		Branch:     fmt.Sprintf("gradle-%s-%s", version, owner),
	}
	newHeadBranch, err := interactor.commitAndPush(ctx, baseBranch, headBranch, commit, files)
	if err != nil {
		return errors.Wrapf(err, "Could not commit and push a branch %s", headBranch)
	}
	log.Debugf(ctx, "Pushed a commit to branch %v", newHeadBranch)

	pull := domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			Repository: base.RepositoryIdentifier,
		},
		Head:  headBranch,
		Base:  baseBranch,
		Title: fmt.Sprintf("Gradle %s", version),
		Body: fmt.Sprintf(`This will upgrade the Gradle wrapper to the latest version %s.

This pull request is sent by @gradleupdate and based on [int128/latest-gradle-wrapper](https://github.com/int128/latest-gradle-wrapper).
`, version),
	}

	pullRequestService := pullRequestService{interactor.PullRequest}
	newPull, err := pullRequestService.createOrUpdatePullRequest(ctx, pull)
	if err != nil {
		return errors.Wrapf(err, "Could not open a pull request %s", pull.String())
	}
	log.Debugf(ctx, "Opened a pull request %v", newPull.PullRequestIdentifier)
	return nil
}

func (interactor *SendPullRequestForUpdate) downloadGradleWrapperFiles(ctx context.Context, id domain.RepositoryIdentifier) ([]domain.File, error) {
	files := make([]domain.File, len(gradleWrapperFiles))
	for i, source := range gradleWrapperFiles {
		file, err := interactor.Repository.GetFile(ctx, id, source.Path)
		if err != nil {
			return nil, errors.Wrapf(err, "Could not get file %s", source.Path)
		}
		files[i] = source
		files[i].Content = file.Content
	}
	return files, nil
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
		newHeadTree, err := interactor.Tree.Create(ctx, head.Repository, baseCommit.Tree, files)
		if err != nil {
			return domain.Branch{}, errors.Wrapf(err, "Could not create a tree on %v", baseCommit.Tree)
		}
		commit.Tree = newHeadTree
		newHeadCommit, err := interactor.Commit.Create(ctx, commit)
		if err != nil {
			return domain.Branch{}, errors.Wrapf(err, "Could not create a commit %v", commit)
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
	newHeadTree, err := interactor.Tree.Create(ctx, head.Repository, baseCommit.Tree, files)
	if err != nil {
		return domain.Branch{}, errors.Wrapf(err, "Could not create a tree on %v", baseCommit.Tree)
	}
	commit.Tree = newHeadTree
	newHeadCommit, err := interactor.Commit.Create(ctx, commit)
	if err != nil {
		return domain.Branch{}, errors.Wrapf(err, "Could not create a commit %v", commit)
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

type pullRequestService struct {
	PullRequest repositories.PullRequest
}

func (service *pullRequestService) createOrUpdatePullRequest(ctx context.Context, pull domain.PullRequest) (domain.PullRequest, error) {
	pulls, err := service.PullRequest.Query(ctx, repositories.PullRequestQuery{
		Head:      pull.Head,
		Base:      pull.Base,
		State:     "open",
		Direction: "desc",
		Sort:      "updated",
		Page:      1,
		PerPage:   1,
	})
	if err != nil {
		return domain.PullRequest{}, errors.Wrapf(err, "Could not find the pull request %s", pull.PullRequestIdentifier.String())
	}
	if len(pulls) > 1 {
		return domain.PullRequest{}, errors.Errorf("Expect single but got %d pull requests", len(pulls))
	}
	if len(pulls) == 1 {
		existent := pulls[0]
		existent.Body = pull.Body
		existent.Title = pull.Title
		updated, err := service.PullRequest.Update(ctx, existent)
		if err != nil {
			return domain.PullRequest{}, errors.Wrapf(err, "Could not update the pull request %s", pull.PullRequestIdentifier.String())
		}
		return updated, err
	}
	created, err := service.PullRequest.Create(ctx, pull)
	if err != nil {
		return domain.PullRequest{}, errors.Wrapf(err, "Could not create a pull request on the repository %s", pull.Repository.String())
	}
	return created, nil
}
