package gateways

import (
	"context"
	"encoding/base64"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

type GitService struct {
	dig.In
	Client *github.Client
}

func (r *GitService) CreateBranch(ctx context.Context, req gateways.PushBranchRequest) (*git.Branch, error) {
	headCommit, err := r.createCommit(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create a commit")
	}
	headRef, _, err := r.Client.Git.CreateRef(ctx,
		req.HeadBranch.Repository.Owner,
		req.HeadBranch.Repository.Name,
		&github.Reference{
			Ref:    github.String(req.HeadBranch.Ref()),
			Object: &github.GitObject{SHA: headCommit.SHA},
		})
	if err != nil {
		return nil, errors.Wrapf(err, "could not create a ref %s in the head repository %s", req.HeadBranch.Ref(), req.HeadBranch.Repository)
	}
	return &git.Branch{
		ID: req.HeadBranch,
		Commit: git.Commit{
			ID: git.CommitID{
				Repository: req.HeadBranch.Repository,
				SHA:        git.CommitSHA(headRef.GetObject().GetSHA()),
			},
			Parents: []git.CommitID{req.BaseBranch.Commit.ID},
			Tree: git.TreeID{
				Repository: req.HeadBranch.Repository,
				SHA:        git.TreeSHA(headCommit.Tree.GetSHA()),
			},
		},
	}, nil
}

func (r *GitService) UpdateForceBranch(ctx context.Context, req gateways.PushBranchRequest) (*git.Branch, error) {
	headCommit, err := r.createCommit(ctx, req)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create a commit")
	}
	headRef, _, err := r.Client.Git.UpdateRef(ctx,
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
	return &git.Branch{
		ID: req.HeadBranch,
		Commit: git.Commit{
			ID: git.CommitID{
				Repository: req.HeadBranch.Repository,
				SHA:        git.CommitSHA(headRef.GetObject().GetSHA()),
			},
			Parents: []git.CommitID{req.BaseBranch.Commit.ID},
			Tree: git.TreeID{
				Repository: req.HeadBranch.Repository,
				SHA:        git.TreeSHA(headCommit.Tree.GetSHA()),
			},
		},
	}, nil
}

func (r *GitService) createCommit(ctx context.Context, req gateways.PushBranchRequest) (*github.Commit, error) {
	headTreeEntries := make([]github.TreeEntry, len(req.CommitFiles))
	for i, file := range req.CommitFiles {
		content := base64.StdEncoding.EncodeToString(file.Content)
		blob, _, err := r.Client.Git.CreateBlob(ctx,
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
	headTree, _, err := r.Client.Git.CreateTree(ctx,
		req.HeadBranch.Repository.Owner,
		req.HeadBranch.Repository.Name,
		req.BaseBranch.Commit.Tree.SHA.String(),
		headTreeEntries)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create a tree in the head repository %s", req.HeadBranch.Repository)
	}
	headCommit, _, err := r.Client.Git.CreateCommit(ctx,
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
