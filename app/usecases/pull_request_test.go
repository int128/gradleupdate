package usecases

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/app/domain"
	"github.com/int128/gradleupdate/app/usecases/mock_repositories"
	"reflect"
	"testing"
)

func TestPullRequestService_createOrUpdatePullRequest_IfNotExist(t *testing.T) {
	ctx := context.Background()
	mc := gomock.NewController(t)
	defer mc.Finish()

	r := mock_repositories.NewMockPullRequest(mc)
	r.EXPECT().Query(ctx, gomock.Any()).Return([]domain.PullRequest{}, nil)
	r.EXPECT().Create(ctx, domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
		},
		Head: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "ho", Repo: "hr"},
			Branch:               "h",
		},
		Base: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
			Branch:               "b",
		},
		Title: "t",
		Body:  "b",
	}).Return(domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
			PullRequestNumber:    3,
		},
		Head: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "ho", Repo: "hr"},
			Branch:               "h",
		},
		Base: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
			Branch:               "b",
		},
		Title: "t",
		Body:  "b",
	}, nil)

	s := pullRequestService{r}
	created, err := s.createOrUpdatePullRequest(ctx, domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
		},
		Head: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "ho", Repo: "hr"},
			Branch:               "h",
		},
		Base: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
			Branch:               "b",
		},
		Title: "t",
		Body:  "b",
	})
	if err != nil {
		t.Fatalf("createOrUpdatePullRequest returned error %s", err)
	}
	if !reflect.DeepEqual(&created, &domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
			PullRequestNumber:    3,
		},
		Head: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "ho", Repo: "hr"},
			Branch:               "h",
		},
		Base: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
			Branch:               "b",
		},
		Title: "t",
		Body:  "b",
	}) {
		t.Errorf("createOrUpdatePullRequest returned wrong value %+v", created)
	}
}

func TestPullRequestService_createOrUpdatePullRequest_IfExists(t *testing.T) {
	ctx := context.Background()
	mc := gomock.NewController(t)
	defer mc.Finish()

	r := mock_repositories.NewMockPullRequest(mc)
	r.EXPECT().Query(ctx, gomock.Any()).Return([]domain.PullRequest{
		{
			PullRequestIdentifier: domain.PullRequestIdentifier{
				RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
				PullRequestNumber:    5,
			},
			Head: domain.BranchIdentifier{
				RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "ho", Repo: "hr"},
				Branch:               "h",
			},
			Base: domain.BranchIdentifier{
				RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
				Branch:               "b",
			},
			Title: "t",
			Body:  "b",
		},
	}, nil)
	r.EXPECT().Update(ctx, domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
			PullRequestNumber:    5,
		},
		Head: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "ho", Repo: "hr"},
			Branch:               "h",
		},
		Base: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
			Branch:               "b",
		},
		Title: "t2",
		Body:  "b2",
	}).Return(domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
			PullRequestNumber:    3,
		},
		Head: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "ho", Repo: "hr"},
			Branch:               "h",
		},
		Base: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
			Branch:               "b",
		},
		Title: "t2",
		Body:  "b2",
	}, nil)

	s := pullRequestService{r}
	updated, err := s.createOrUpdatePullRequest(ctx, domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
		},
		Head: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "ho", Repo: "hr"},
			Branch:               "h",
		},
		Base: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
			Branch:               "b",
		},
		Title: "t",
		Body:  "b",
	})
	if err != nil {
		t.Fatalf("createOrUpdatePullRequest returned error %s", err)
	}
	if !reflect.DeepEqual(&updated, &domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
			PullRequestNumber:    3,
		},
		Head: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "ho", Repo: "hr"},
			Branch:               "h",
		},
		Base: domain.BranchIdentifier{
			RepositoryIdentifier: domain.RepositoryIdentifier{Owner: "bo", Repo: "br"},
			Branch:               "b",
		},
		Title: "t2",
		Body:  "b2",
	}) {
		t.Errorf("createOrUpdatePullRequest returned wrong value %+v", updated)
	}
}
