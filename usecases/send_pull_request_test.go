package usecases

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/gateways/mock_gateways"
)

func TestPullRequestService_createOrUpdatePullRequest_IfNotExist(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_gateways.NewMockPullRequestRepository(ctrl)
	r.EXPECT().Query(ctx, gomock.Any()).Return([]domain.PullRequest{}, nil)
	r.EXPECT().Create(ctx, domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
		},
		HeadBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "ho", Name: "hr"},
			Name:       "h",
		},
		BaseBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
			Name:       "b",
		},
		Title: "t",
		Body:  "b",
	}).Return(&domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			Repository:        domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
			PullRequestNumber: 3,
		},
		HeadBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "ho", Name: "hr"},
			Name:       "h",
		},
		BaseBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
			Name:       "b",
		},
		Title: "t",
		Body:  "b",
	}, nil)

	s := pullRequestService{r}
	created, err := s.createOrUpdatePullRequest(ctx, domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
		},
		HeadBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "ho", Name: "hr"},
			Name:       "h",
		},
		BaseBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
			Name:       "b",
		},
		Title: "t",
		Body:  "b",
	})
	if err != nil {
		t.Fatalf("createOrUpdatePullRequest returned error %s", err)
	}
	if !reflect.DeepEqual(created, &domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			Repository:        domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
			PullRequestNumber: 3,
		},
		HeadBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "ho", Name: "hr"},
			Name:       "h",
		},
		BaseBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
			Name:       "b",
		},
		Title: "t",
		Body:  "b",
	}) {
		t.Errorf("createOrUpdatePullRequest returned wrong value %+v", created)
	}
}

func TestPullRequestService_createOrUpdatePullRequest_IfExists(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := mock_gateways.NewMockPullRequestRepository(ctrl)
	r.EXPECT().Query(ctx, gomock.Any()).Return([]domain.PullRequest{
		{
			PullRequestIdentifier: domain.PullRequestIdentifier{
				Repository:        domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
				PullRequestNumber: 5,
			},
			HeadBranch: domain.BranchIdentifier{
				Repository: domain.RepositoryIdentifier{Owner: "ho", Name: "hr"},
				Name:       "h",
			},
			BaseBranch: domain.BranchIdentifier{
				Repository: domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
				Name:       "b",
			},
			Title: "t",
			Body:  "b",
		},
	}, nil)
	r.EXPECT().Update(ctx, domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			Repository:        domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
			PullRequestNumber: 5,
		},
		HeadBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "ho", Name: "hr"},
			Name:       "h",
		},
		BaseBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
			Name:       "b",
		},
		Title: "t",
		Body:  "b",
	}).Return(&domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			Repository:        domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
			PullRequestNumber: 3,
		},
		HeadBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "ho", Name: "hr"},
			Name:       "h",
		},
		BaseBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
			Name:       "b",
		},
		Title: "t2",
		Body:  "b2",
	}, nil)

	s := pullRequestService{r}
	updated, err := s.createOrUpdatePullRequest(ctx, domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
		},
		HeadBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "ho", Name: "hr"},
			Name:       "h",
		},
		BaseBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
			Name:       "b",
		},
		Title: "t",
		Body:  "b",
	})
	if err != nil {
		t.Fatalf("createOrUpdatePullRequest returned error %s", err)
	}
	if !reflect.DeepEqual(updated, &domain.PullRequest{
		PullRequestIdentifier: domain.PullRequestIdentifier{
			Repository:        domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
			PullRequestNumber: 3,
		},
		HeadBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "ho", Name: "hr"},
			Name:       "h",
		},
		BaseBranch: domain.BranchIdentifier{
			Repository: domain.RepositoryIdentifier{Owner: "bo", Name: "br"},
			Name:       "b",
		},
		Title: "t2",
		Body:  "b2",
	}) {
		t.Errorf("createOrUpdatePullRequest returned wrong value %+v", updated)
	}
}
