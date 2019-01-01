package gateways

import (
	"context"
	"encoding/base64"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/infrastructure/interfaces"
	"github.com/pkg/errors"
)

type GitService struct {
	GitHubClientFactory infrastructure.GitHubClientFactory
}

func (r *GitService) CreateBranch(ctx context.Context, req gateways.PushBranchRequest) (*domain.Branch, error) {
	client := r.GitHubClientFactory.New(ctx)
	headCommit, err := createCommit(ctx, client, req)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create a commit")
	}
	headRef, _, err := client.Git.CreateRef(ctx,
		req.HeadBranch.Repository.Owner,
		req.HeadBranch.Repository.Name,
		&github.Reference{
			Ref:    github.String(req.HeadBranch.Ref()),
			Object: &github.GitObject{SHA: headCommit.SHA},
		})
	if err != nil {
		return nil, errors.Wrapf(err, "could not create a ref %s in the head repository %s", req.HeadBranch.Ref(), req.HeadBranch.Repository)
	}
	return &domain.Branch{
		ID: req.HeadBranch,
		Commit: domain.Commit{
			ID: domain.CommitID{
				Repository: req.HeadBranch.Repository,
				SHA:        domain.CommitSHA(headRef.GetObject().GetSHA()),
			},
			Parents: []domain.CommitID{req.BaseBranch.Commit.ID},
			Tree: domain.TreeID{
				Repository: req.HeadBranch.Repository,
				SHA:        domain.TreeSHA(headCommit.Tree.GetSHA()),
			},
		},
	}, nil
}

func (r *GitService) UpdateForceBranch(ctx context.Context, req gateways.PushBranchRequest) (*domain.Branch, error) {
	client := r.GitHubClientFactory.New(ctx)
	headCommit, err := createCommit(ctx, client, req)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create a commit")
	}
	headRef, _, err := client.Git.UpdateRef(ctx,
		req.HeadBranch.Repository.Owner,
		req.HeadBranch.Repository.Name,
		&github.Reference{
			Ref:    github.String(req.HeadBranch.Ref()),
			Object: &github.GitObject{SHA: headCommit.SHA},
		},
		true)
	if err != nil {
		return nil, errors.Wrapf(err, "could not update the ref %s in the head repository %s", req.HeadBranch.Ref(), req.HeadBranch.Repository)
	}
	return &domain.Branch{
		ID: req.HeadBranch,
		Commit: domain.Commit{
			ID: domain.CommitID{
				Repository: req.HeadBranch.Repository,
				SHA:        domain.CommitSHA(headRef.GetObject().GetSHA()),
			},
			Parents: []domain.CommitID{req.BaseBranch.Commit.ID},
			Tree: domain.TreeID{
				Repository: req.HeadBranch.Repository,
				SHA:        domain.TreeSHA(headCommit.Tree.GetSHA()),
			},
		},
	}, nil
}

func createCommit(ctx context.Context, client *github.Client, req gateways.PushBranchRequest) (*github.Commit, error) {
	headTreeEntries := make([]github.TreeEntry, len(req.CommitFiles))
	for i, file := range req.CommitFiles {
		content := base64.StdEncoding.EncodeToString(file.Content)
		blob, _, err := client.Git.CreateBlob(ctx,
			req.HeadBranch.Repository.Owner,
			req.HeadBranch.Repository.Name,
			&github.Blob{
				Content:  github.String(content),
				Encoding: github.String("base64"),
			})
		if err != nil {
			return nil, errors.Wrapf(err, "could not create a blob in the head repository %s", req.HeadBranch.Repository)
		}
		headTreeEntries[i] = github.TreeEntry{
			Path: github.String(file.Path),
			Mode: github.String("100644"),
			SHA:  blob.SHA,
		}
	}
	headTree, _, err := client.Git.CreateTree(ctx,
		req.HeadBranch.Repository.Owner,
		req.HeadBranch.Repository.Name,
		req.BaseBranch.Commit.Tree.SHA.String(),
		headTreeEntries)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create a tree in the head repository %s", req.HeadBranch.Repository)
	}
	headCommit, _, err := client.Git.CreateCommit(ctx,
		req.HeadBranch.Repository.Owner,
		req.HeadBranch.Repository.Name,
		&github.Commit{
			Parents: []github.Commit{{
				SHA: github.String(req.BaseBranch.Commit.ID.SHA.String()),
			}},
			Message: github.String(req.CommitMessage),
			Tree:    headTree,
		})
	if err != nil {
		return nil, errors.Wrapf(err, "could not create a commit in the head repository %s", req.HeadBranch.Repository)
	}
	return headCommit, nil
}
