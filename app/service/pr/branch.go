package pr

import (
	"context"

	"github.com/google/go-github/v18/github"
	"github.com/pkg/errors"
)

// CreateOrUpdateBranch creates a head branch with a commit on the base if the head does not exist.
// Otherwise it updates the head branch if it is out-of-date, that is,
// parent of the head does not equal to the base.
func CreateOrUpdateBranch(ctx context.Context, c *github.Client, base, head Branch, commit Commit) error {
	headRef, resp, err := c.Git.GetRef(ctx, head.Owner, head.Repo, "refs/heads/"+head.Branch)
	if resp != nil && resp.StatusCode == 404 {
		if _, err := createBranch(ctx, c, base, head, commit); err != nil {
			return errors.Wrapf(err, "Could not create a branch")
		}
		return nil
	}
	if err != nil {
		return errors.Wrapf(err, "Could not get the ref of head branch")
	}
	if _, err := updateBranchIfHeadIsOutOfDate(ctx, c, base, head, headRef, commit); err != nil {
		return errors.Wrapf(err, "Could not update the branch")
	}
	return nil
}

func createBranch(ctx context.Context, c *github.Client, base, head Branch, commit Commit) (*github.Reference, error) {
	baseRef, _, err := c.Git.GetRef(ctx, base.Owner, base.Repo, "refs/heads/"+base.Branch)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get the ref of base branch %s", base)
	}
	baseCommitSHA := baseRef.GetObject().GetSHA()
	baseCommit, _, err := c.Git.GetCommit(ctx, base.Owner, base.Repo, baseCommitSHA)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get the commit of base branch %s sha %s", base, baseCommitSHA)
	}
	baseTreeSHA := baseCommit.GetTree().GetSHA()
	newHeadCommit, err := createCommit(ctx, c, head.Repository, baseCommitSHA, baseTreeSHA, commit)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create a commit on head repository %s", head.Repository)
	}
	ref, _, err := c.Git.CreateRef(ctx, head.Owner, head.Repo, &github.Reference{
		Ref:    github.String("refs/heads/" + head.Branch),
		Object: &github.GitObject{SHA: newHeadCommit.SHA},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create a ref for head branch %s", head)
	}
	return ref, nil
}

func updateBranchIfHeadIsOutOfDate(ctx context.Context, c *github.Client, base, head Branch, headRef *github.Reference, commit Commit) (*github.Reference, error) {
	baseRef, _, err := c.Git.GetRef(ctx, base.Owner, base.Repo, "refs/heads/"+base.Branch)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get the ref of base branch %s", base)
	}
	baseCommitSHA := baseRef.GetObject().GetSHA()

	headCommitSHA := headRef.GetObject().GetSHA()
	headCommit, _, err := c.Git.GetCommit(ctx, head.Owner, head.Repo, headCommitSHA)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get the commit of head branch %s sha %s", head, headCommitSHA)
	}
	headIsOutOfDate := determineIfHeadIsOutOfDate(headCommit, baseCommitSHA)
	if !headIsOutOfDate {
		return headRef, nil
	}

	baseCommit, _, err := c.Git.GetCommit(ctx, base.Owner, base.Repo, baseCommitSHA)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not get the commit of base branch %s sha %s", base, baseCommitSHA)
	}
	baseTreeSHA := baseCommit.GetTree().GetSHA()
	newHeadCommit, err := createCommit(ctx, c, head.Repository, baseCommitSHA, baseTreeSHA, commit)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create a commit on head repository %s", head.Repository)
	}
	ref, _, err := c.Git.UpdateRef(ctx, head.Owner, head.Repo, &github.Reference{
		Ref:    github.String("refs/heads/" + head.Branch),
		Object: &github.GitObject{SHA: newHeadCommit.SHA},
	}, true)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not update the ref for head branch %s", head)
	}
	return ref, nil
}

func determineIfHeadIsOutOfDate(headCommit *github.Commit, baseCommitSHA string) bool {
	if len(headCommit.Parents) != 1 {
		return true
	}
	return headCommit.Parents[0].GetSHA() != baseCommitSHA
}
