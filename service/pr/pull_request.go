package pr

import (
	"context"

	"github.com/google/go-github/v18/github"
	"github.com/pkg/errors"
)

// CreateOrUpdatePullRequest creates a pull request if it does not exist,
// or updates the pull request updated most recently.
func CreateOrUpdatePullRequest(ctx context.Context, c *github.Client, pr PullRequest) (*github.PullRequest, error) {
	pulls, _, err := c.PullRequests.List(ctx, pr.Base.Owner, pr.Base.Repo, &github.PullRequestListOptions{
		Base:        pr.Base.Branch,
		Head:        pr.Head.Owner + ":" + pr.Head.Branch,
		State:       "open",
		Direction:   "desc",
		Sort:        "updated",
		ListOptions: github.ListOptions{PerPage: 1, Page: 1},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not find the pull request %s", pr)
	}
	if len(pulls) > 1 {
		return nil, errors.Errorf("Expect single but got %d pull requests", len(pulls))
	}
	if len(pulls) == 1 {
		pull, _, err := c.PullRequests.Edit(ctx, pr.Base.Owner, pr.Base.Repo, pulls[0].GetNumber(), &github.PullRequest{
			Title: github.String(pr.Title),
			Body:  github.String(pr.Body),
		})
		if err != nil {
			return nil, errors.Wrapf(err, "Could not update the pull request %s", pr)
		}
		return pull, nil
	}
	pull, _, err := c.PullRequests.Create(ctx, pr.Base.Owner, pr.Base.Repo, &github.NewPullRequest{
		Base:  github.String(pr.Base.Branch),
		Head:  github.String(pr.Head.Owner + ":" + pr.Head.Branch),
		Title: github.String(pr.Title),
		Body:  github.String(pr.Body),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create a pull request %s", pr)
	}
	return pull, nil
}
