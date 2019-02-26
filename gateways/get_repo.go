package gateways

import (
	"context"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/pkg/errors"
	"github.com/shurcooL/githubv4"
	"go.uber.org/dig"
)

type GetRepositoryQueryIn struct {
	dig.In
	Client *githubv4.Client
}

type GetRepositoryQuery struct {
	GetRepositoryQueryIn
	noSuchEntityErrorCauser
}

func (r *GetRepositoryQuery) Do(ctx context.Context, in gateways.GetRepositoryQueryIn) (*gateways.GetRepositoryQueryOut, error) {
	type blobObject struct {
		Blob struct {
			Text string
		} `graphql:"... on Blob"`
	}
	var q *struct {
		Repository struct {
			Name  string
			Owner struct {
				Login     string
				AvatarUrl string
			}
			Url         string
			Description string

			// a pull request associated with the head branch
			PullRequests struct {
				Nodes []struct {
					Url string
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

	var out gateways.GetRepositoryQueryOut
	out.Repository = git.Repository{
		ID: git.RepositoryID{
			Owner: q.Repository.Owner.Login,
			Name:  q.Repository.Name,
		},
		Description: q.Repository.Description,
		AvatarURL:   q.Repository.Owner.AvatarUrl,
		URL:         q.Repository.Url,
		//TODO: DefaultBranch is omitted because repo page does not need it
	}
	if len(q.Repository.PullRequests.Nodes) == 1 {
		pull := q.Repository.PullRequests.Nodes[0]
		out.PullRequestURL = git.PullRequestURL(pull.Url)
	}
	if q.Repository.ReadmeMd != nil {
		out.Readme = git.FileContent(q.Repository.ReadmeMd.Blob.Text)
	}
	if q.Repository.GradleWrapperProperties != nil {
		out.GradleWrapperProperties = git.FileContent(q.Repository.GradleWrapperProperties.Blob.Text)
	}
	return &out, nil
}
