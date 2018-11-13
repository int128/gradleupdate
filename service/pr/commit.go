package pr

import (
	"context"

	"github.com/google/go-github/v18/github"
	"github.com/pkg/errors"
)

func createCommit(ctx context.Context, c *github.Client, r Repository, parentCommitSHA, parentTreeSHA string, commit Commit) (*github.Commit, error) {
	treeEntries, err := createTreeEntries(ctx, c, r, commit.Files)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create tree entries")
	}
	tree, _, err := c.Git.CreateTree(ctx, r.Owner, r.Repo, parentTreeSHA, treeEntries)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create a tree on repository %s", r)
	}
	created, _, err := c.Git.CreateCommit(ctx, r.Owner, r.Repo, &github.Commit{
		Message: github.String(commit.Message),
		Tree:    tree,
		Parents: []github.Commit{github.Commit{SHA: github.String(parentCommitSHA)}},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Could not create a commit on repository %s", r)
	}
	return created, nil
}

func createTreeEntries(ctx context.Context, c *github.Client, r Repository, files []File) ([]github.TreeEntry, error) {
	e := make([]github.TreeEntry, len(files))
	for i, file := range files {
		blob, _, err := c.Git.CreateBlob(ctx, r.Owner, r.Repo, &github.Blob{
			Encoding: github.String("base64"),
			Content:  github.String(file.EncodedContent),
		})
		if err != nil {
			return nil, errors.Wrapf(err, "Could not create a blob for repository %s file %s", r, file.Path)
		}
		e[i] = github.TreeEntry{
			Path: github.String(file.Path),
			Mode: github.String(file.Mode),
			SHA:  blob.SHA,
		}
	}
	return e, nil
}
