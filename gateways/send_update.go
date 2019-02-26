package gateways

import (
	"context"
	"encoding/base64"

	"github.com/google/go-github/v24/github"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
	"go.uber.org/dig"
)

type SendUpdateQueryIn struct {
	dig.In
	Client   *githubv4.Client
	ClientV3 *github.Client
}

// SendUpdateQuery provides GitHub access for the SendUpdate usecase.
type SendUpdateQuery struct {
	SendUpdateQueryIn
	noSuchEntityErrorCauser
}

func (r *SendUpdateQuery) Get(ctx context.Context, in gateways.SendUpdateQueryIn) (*gateways.SendUpdateQueryOut, error) {
	type blobObject struct {
		Blob struct {
			Text string
		} `graphql:"... on Blob"`
	}
	var q *struct {
		Repository struct {
			Name  string
			Owner struct{ Login string }

			// default branch (name, commit SHA and tree SHA)
			DefaultBranchRef struct {
				Name   string
				Target struct {
					Commit struct {
						Oid  string
						Tree struct {
							Oid string
						}
					} `graphql:"... on Commit"`
				}
			}

			// a pull request associated with the head branch
			PullRequests struct {
				Nodes []struct {
					HeadRef struct {
						Name       string
						Repository struct {
							Name  string
							Owner struct{ Login string }
						}
						Target struct {
							Commit struct {
								// the parent of the head ref
								Parents struct {
									TotalCount int
									Nodes      []struct {
										Oid string
									}
								} `graphql:"parents(first: 1)"`
							} `graphql:"... on Commit"`
						}
					}
				}
			} `graphql:"pullRequests(first: 1, headRefName: $headRefName)"`

			// files
			ReadmeMd                *blobObject `graphql:"readmeMd: object(expression: $readmeMd)"`
			GradleWrapperProperties *blobObject `graphql:"gradleWrapperProperties: object(expression: $gradleWrapperProperties)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	v := map[string]interface{}{
		"owner":                   githubv4.String(in.Repository.Owner),
		"repo":                    githubv4.String(in.Repository.Name),
		"headRefName":             githubv4.String(in.HeadBranchName),
		"readmeMd":                githubv4.String("HEAD:README.md"),
		"gradleWrapperProperties": githubv4.String("HEAD:" + gradle.WrapperPropertiesPath),
	}
	if err := r.Client.Query(ctx, &q, v); err != nil {
		// githubv4 does not provide errors interface for now,
		// so we check the pointer is nil on not found error.
		// See https://github.com/shurcooL/githubv4/issues/41
		if q != nil {
			return nil, errors.Wrapf(&noSuchEntityError{err}, "no such repository %s", in.Repository)
		}
		return nil, errors.Wrapf(err, "GitHub API error")
	}

	var out gateways.SendUpdateQueryOut
	out.BaseRepository = git.RepositoryID{
		Owner: q.Repository.Owner.Login,
		Name:  q.Repository.Name,
	}
	out.BaseBranch = git.BranchID{
		Repository: git.RepositoryID{
			Owner: q.Repository.Owner.Login,
			Name:  q.Repository.Name,
		},
		Name: git.BranchName(q.Repository.DefaultBranchRef.Name),
	}
	out.BaseCommitSHA = git.CommitSHA(q.Repository.DefaultBranchRef.Target.Commit.Oid)
	out.BaseTreeSHA = git.TreeSHA(q.Repository.DefaultBranchRef.Target.Commit.Tree.Oid)
	if len(q.Repository.PullRequests.Nodes) == 1 {
		pull := q.Repository.PullRequests.Nodes[0]
		out.HeadBranch = &git.BranchID{
			Repository: git.RepositoryID{
				Owner: pull.HeadRef.Repository.Owner.Login,
				Name:  pull.HeadRef.Repository.Name,
			},
			Name: git.BranchName(pull.HeadRef.Name),
		}
		parents := pull.HeadRef.Target.Commit.Parents
		if parents.TotalCount == 1 && len(parents.Nodes) == 1 {
			out.HeadParentCommitSHA = git.CommitSHA(parents.Nodes[0].Oid)
		}
	}
	if q.Repository.ReadmeMd != nil {
		out.Readme = git.FileContent(q.Repository.ReadmeMd.Blob.Text)
	}
	if q.Repository.GradleWrapperProperties != nil {
		out.GradleWrapperProperties = git.FileContent(q.Repository.GradleWrapperProperties.Blob.Text)
	}
	return &out, nil
}

func (r *SendUpdateQuery) ForkRepository(ctx context.Context, id git.RepositoryID) (*git.RepositoryID, error) {
	fork, _, err := r.ClientV3.Repositories.CreateFork(ctx, id.Owner, id.Name, &github.RepositoryCreateForkOptions{})
	if err != nil {
		if _, ok := err.(*github.AcceptedError); ok {
			// Fork in progress
		} else {
			return nil, errors.Wrapf(err, "error from GitHub API")
		}
	}
	return &git.RepositoryID{
		Owner: fork.GetOwner().GetLogin(),
		Name:  fork.GetName(),
	}, nil
}

func (r *SendUpdateQuery) CreateBranch(ctx context.Context, b gateways.NewBranch) error {
	commit, err := r.createCommit(ctx, b)
	if err != nil {
		return errors.Wrapf(err, "error while creating a commit")
	}
	ref := &github.Reference{
		Ref:    github.String(b.Branch.Name.Ref()),
		Object: &github.GitObject{SHA: commit.SHA},
	}
	if _, _, err := r.ClientV3.Git.CreateRef(ctx, b.Branch.Repository.Owner, b.Branch.Repository.Name, ref); err != nil {
		return errors.Wrapf(err, "error while creating a branch %s", b.Branch)
	}
	return nil
}

func (r *SendUpdateQuery) UpdateBranch(ctx context.Context, b gateways.NewBranch, force bool) error {
	commit, err := r.createCommit(ctx, b)
	if err != nil {
		return errors.Wrapf(err, "error while creating a commit")
	}
	refIn := github.Reference{
		Ref:    github.String(b.Branch.Name.Ref()),
		Object: &github.GitObject{SHA: commit.SHA},
	}
	if _, _, err := r.ClientV3.Git.UpdateRef(ctx, b.Branch.Repository.Owner, b.Branch.Repository.Name, &refIn, force); err != nil {
		return errors.Wrapf(err, "error while updating the branch %s", b.Branch)
	}
	return nil
}

func (r *SendUpdateQuery) createCommit(ctx context.Context, b gateways.NewBranch) (*github.Commit, error) {
	treeEntries := make([]github.TreeEntry, len(b.CommitFiles))
	for i, file := range b.CommitFiles {
		blobIn := github.Blob{
			Content:  github.String(base64.StdEncoding.EncodeToString(file.Content)),
			Encoding: github.String("base64"),
		}
		blob, _, err := r.ClientV3.Git.CreateBlob(ctx, b.Branch.Repository.Owner, b.Branch.Repository.Name, &blobIn)
		if err != nil {
			return nil, errors.Wrapf(err, "error while creating a blob in the repository %s", b.Branch.Repository)
		}
		treeEntries[i] = github.TreeEntry{
			Path: github.String(file.Path),
			Mode: github.String("100644"),
			SHA:  blob.SHA,
		}
	}

	tree, _, err := r.ClientV3.Git.CreateTree(ctx, b.Branch.Repository.Owner, b.Branch.Repository.Name, b.ParentTreeSHA.String(), treeEntries)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating a tree in the repository %s", b.Branch.Repository)
	}

	commitIn := github.Commit{
		Tree:    tree,
		Parents: []github.Commit{{SHA: github.String(b.ParentCommitSHA.String())}},
		Message: github.String(b.CommitMessage),
	}
	commit, _, err := r.ClientV3.Git.CreateCommit(ctx, b.Branch.Repository.Owner, b.Branch.Repository.Name, &commitIn)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating a commit in the repository %s", b.Branch.Repository)
	}
	return commit, nil
}
