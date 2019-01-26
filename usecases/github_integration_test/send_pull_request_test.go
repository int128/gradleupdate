package usecases_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways"
	"github.com/int128/gradleupdate/usecases"
	interfaces "github.com/int128/gradleupdate/usecases/interfaces"
)

func TestSendPullRequestRequest_Do(t *testing.T) {
	ctx := context.Background()
	client := newGitHubClient(t)
	sendPullRequest := usecases.SendPullRequest{
		RepositoryRepository:  &gateways.RepositoryRepository{Client: client},
		PullRequestRepository: &gateways.PullRequestRepository{Client: client},
		GitService:            &gateways.GitService{Client: client},
	}

	// Delete the head branch if it exists
	if _, err := client.Git.DeleteRef(ctx, forkedRepository.Owner, forkedRepository.Name, "refs/heads/example"); err != nil {
		if resp, ok := err.(*github.ErrorResponse); ok {
			if resp.Message == "Reference does not exist" {
				// did not exist
			} else {
				t.Fatalf("could not delete the head branch: %s", err)
			}
		} else {
			t.Fatalf("could not delete the head branch: %s", err)
		}
	}

	req := interfaces.SendPullRequestRequest{
		Base:           sandboxRepository,
		HeadBranchName: "example",
		CommitMessage:  "Example Commit",
		CommitFiles: []domain.File{
			{
				Path:    "foo/bar",
				Content: domain.FileContent("baz"),
			},
		},
		Title: "Example",
		Body:  "This is an example pull request.",
	}
	if err := sendPullRequest.Do(ctx, req); err != nil {
		t.Fatalf("could not do the usecase: %s", err)
	}
	assertHeadBranchExists(t, ctx, client, req)
	assertPullRequestExists(t, ctx, client, req)

	if err := sendPullRequest.Do(ctx, req); err != nil {
		t.Fatalf("could not do the usecase: %s", err)
	}
	assertHeadBranchExists(t, ctx, client, req)
	assertPullRequestExists(t, ctx, client, req)
}

func assertPullRequestExists(t *testing.T, ctx context.Context, client *github.Client, req interfaces.SendPullRequestRequest) {
	t.Helper()
	pulls, _, err := client.PullRequests.List(ctx, sandboxRepository.Owner, sandboxRepository.Name,
		&github.PullRequestListOptions{
			Base:  "master",
			Head:  forkedRepository.Owner + ":example",
			State: "open",
		})
	if err != nil {
		t.Fatalf("could not find the pull request: %s", err)
	}
	if len(pulls) != 1 {
		t.Fatalf("pulls wants 1 but %+v", pulls)
	}
	pull := pulls[0]
	if pull.GetTitle() != req.Title {
		t.Errorf("Title wants %s but %s", req.Title, pull.GetTitle())
	}
	if pull.GetBody() != req.Body {
		t.Errorf("Body wants %s but %s", req.Body, pull.GetBody())
	}
}

func assertHeadBranchExists(t *testing.T, ctx context.Context, client *github.Client, req interfaces.SendPullRequestRequest) {
	t.Helper()
	for _, file := range req.CommitFiles {
		fc, _, _, err := client.Repositories.GetContents(ctx, forkedRepository.Owner, forkedRepository.Name,
			file.Path, &github.RepositoryContentGetOptions{Ref: req.HeadBranchName})
		if err != nil {
			t.Fatalf("could not find %s: %s", file.Path, err)
		}
		content, err := fc.GetContent()
		if err != nil {
			t.Fatalf("could not decode content: %s", err)
		}
		if content != file.Content.String() {
			t.Errorf("content wants baz but %s", content)
		}
	}
}
