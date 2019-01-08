package github_integration_test

import (
	"context"
	"testing"

	"github.com/int128/gradleupdate/domain"
	impl "github.com/int128/gradleupdate/gateways"
	"github.com/int128/gradleupdate/gateways/interfaces"
)

func TestNewPullRequest(t *testing.T) {
	client := newGitHubClient(t)
	gitHubClientFactory := newGitHubClientFactory(t, client)
	repositoryRepository := impl.RepositoryRepository{GitHubClientFactory: gitHubClientFactory}
	pullRequestRepository := impl.PullRequestRepository{GitHubClientFactory: gitHubClientFactory}
	gitService := impl.GitService{GitHubClientFactory: gitHubClientFactory}
	ctx := context.Background()

	// Setup
	t.Logf("Fork the sandbox repository %s", sandboxRepository)
	fork, err := repositoryRepository.Fork(ctx, sandboxRepository)
	if err != nil {
		t.Fatalf("could not fork the sandbox repository: %s", err)
	}

	t.Logf("Get the master branch of %s", sandboxRepository)
	base, err := repositoryRepository.Get(ctx, sandboxRepository)
	if err != nil {
		t.Fatalf("could not get the base repository: %s", err)
	}
	baseBranch, err := repositoryRepository.GetBranch(ctx, base.DefaultBranch)
	if err != nil {
		t.Fatalf("could not get the master branch: %s", err)
	}
	if baseBranch.ID != base.DefaultBranch {
		t.Errorf("ID wants %s but %s", base.DefaultBranch, baseBranch.ID)
	}

	headBranchID := domain.BranchID{
		Repository: fork.ID,
		Name:       "branch1",
	}
	t.Logf("Delete the head branch %s if it exists", headBranchID)
	if resp, err := client.Git.DeleteRef(ctx, fork.ID.Owner, fork.ID.Name, headBranchID.Ref()); err != nil {
		if resp.StatusCode != 404 {
			t.Fatalf("could not delete the ref %s", headBranchID.Ref())
		}
	}

	// Create a branch
	t.Logf("Create a head branch %s based on %s", headBranchID, baseBranch)
	createBranchRequest := gateways.PushBranchRequest{
		BaseBranch:    *baseBranch,
		HeadBranch:    headBranchID,
		CommitMessage: "Test commit",
		CommitFiles: []domain.File{
			{
				Path:    "foo/bar",
				Content: domain.FileContent("baz"),
			},
		},
	}
	branch, err := gitService.CreateBranch(ctx, createBranchRequest)
	if err != nil {
		t.Fatalf("could not create a branch: %s", err)
	}
	if branch.ID != headBranchID {
		t.Errorf("branch.Name wants %s but %s", headBranchID, branch.ID)
	}
	if branch.Commit.ID.SHA == "" {
		t.Errorf("branch.CommitSHA wants non-empty but empty")
	}

	// Create a pull request
	t.Logf("Create a pull request on %s", sandboxRepository)
	createdPull, err := pullRequestRepository.Create(ctx, domain.PullRequest{
		ID:         domain.PullRequestID{Repository: sandboxRepository},
		BaseBranch: baseBranch.ID,
		HeadBranch: headBranchID,
		Title:      "Example",
		Body:       "This is an example pull request.",
	})
	if err != nil {
		t.Fatalf("could not create a pull request: %s", err)
	}
	if createdPull.BaseBranch != baseBranch.ID {
		t.Errorf("BaseBranch wants %s but %s", baseBranch.ID, createdPull.BaseBranch)
	}
	if createdPull.HeadBranch != headBranchID {
		t.Errorf("HeadBranch wants %s but %s", headBranchID, createdPull.HeadBranch)
	}

	t.Logf("Find the created pull request from %s", sandboxRepository)
	foundPull, err := pullRequestRepository.FindByBranch(ctx, baseBranch.ID, headBranchID)
	if err != nil {
		t.Fatalf("could not find the pull request: %s", err)
	}
	if foundPull == nil {
		t.Fatalf("FindByBranch returned nil")
	}
	if foundPull.ID != createdPull.ID {
		t.Errorf("ID wants %s but %s", createdPull.ID, foundPull.ID)
	}

	// Update the branch
	t.Logf("Update head branch %s based on %s", headBranchID, baseBranch)
	updateBranchRequest := gateways.PushBranchRequest{
		BaseBranch:    *baseBranch,
		HeadBranch:    headBranchID,
		CommitMessage: "Test commit",
		CommitFiles: []domain.File{
			{
				Path:    "foo/bar2",
				Content: domain.FileContent("baz"),
			},
		},
	}
	updatedBranch, err := gitService.UpdateForceBranch(ctx, updateBranchRequest)
	if err != nil {
		t.Fatalf("could not update the branch: %s", err)
	}
	if updatedBranch.ID != headBranchID {
		t.Errorf("branch.Name wants %s but %s", headBranchID, updatedBranch.ID)
	}
	if updatedBranch.Commit.ID.SHA == "" {
		t.Errorf("branch.CommitSHA wants non-empty but empty")
	}
}
