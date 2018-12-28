package gateways

import (
	"context"
	"encoding/base64"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/gateways"
	"github.com/pkg/errors"
)

type GitService struct {
	GitHubClientFactory GitHubClientFactory
}

func (r *GitService) ForkBranch(ctx context.Context, req gateways.ForkBranchRequest) (*domain.Branch, error) {
	client := r.GitHubClientFactory.New(ctx)

	head, _, err := client.Repositories.CreateFork(ctx, req.Base.Repository.Owner, req.Base.Repository.Name, nil)
	if err != nil {
		if _, ok := err.(*github.AcceptedError); !ok {
			return nil, errors.Wrapf(err, "could not fork the base repository %s", req.Base.Repository)
		}
	}
	//TODO: wait until fork is completed

	baseRef, _, err := client.Git.GetRef(ctx,
		req.Base.Repository.Owner,
		req.Base.Repository.Name,
		"refs/heads/"+req.Base.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get the base branch %s", req.Base)
	}
	baseCommit, _, err := client.Git.GetCommit(ctx,
		req.Base.Repository.Owner,
		req.Base.Repository.Name,
		baseRef.GetObject().GetSHA())
	if err != nil {
		return nil, errors.Wrapf(err, "could not get the commit %s of base branch %s", baseRef.GetObject().GetSHA(), req.Base)
	}

	headTreeEntries := make([]github.TreeEntry, len(req.Files))
	for i, file := range req.Files {
		content := base64.StdEncoding.EncodeToString(file.Content)
		blob, _, err := client.Git.CreateBlob(ctx,
			head.GetOwner().GetLogin(),
			head.GetName(),
			&github.Blob{
				Content:  github.String(content),
				Encoding: github.String("base64"),
			})
		if err != nil {
			return nil, errors.Wrapf(err, "could not create a blob in the head repository %s", head.GetFullName())
		}
		headTreeEntries[i] = github.TreeEntry{
			Path: github.String(file.Path),
			Mode: github.String("100644"),
			SHA:  blob.SHA,
		}
	}
	headTree, _, err := client.Git.CreateTree(ctx,
		head.GetOwner().GetLogin(),
		head.GetName(),
		baseCommit.GetTree().GetSHA(),
		headTreeEntries)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create a tree in the head repository %s", head.GetFullName())
	}
	headCommit, _, err := client.Git.CreateCommit(ctx,
		head.GetOwner().GetLogin(),
		head.GetName(),
		&github.Commit{
			Parents: []github.Commit{{SHA: baseCommit.SHA}},
			Message: github.String(req.CommitMessage),
			Tree:    headTree,
		})
	if err != nil {
		return nil, errors.Wrapf(err, "could not create a commit in the head repository %s", head.GetFullName())
	}

	headRef, _, err := client.Git.CreateRef(ctx,
		head.GetOwner().GetLogin(),
		head.GetName(),
		&github.Reference{
			Ref:    github.String("refs/heads/" + req.HeadBranchName),
			Object: &github.GitObject{SHA: headCommit.SHA},
		})
	if err != nil {
		return nil, errors.Wrapf(err, "could not create a ref %s in the head repository %s", "refs/heads/"+req.HeadBranchName, head.GetFullName())
	}
	return &domain.Branch{
		BranchIdentifier: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{
				Owner: head.GetOwner().GetLogin(),
				Name:  head.GetName(),
			},
			Name: req.HeadBranchName,
		},
		CommitSHA: headRef.GetObject().GetSHA(),
	}, nil
}
