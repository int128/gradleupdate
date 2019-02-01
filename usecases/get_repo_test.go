package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/testdata"
	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
)

func TestGetRepository_Do(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}
	readme := git.FileContent("/owner/repo/status.svg")

	for _, c := range []struct {
		name    string
		content git.FileContent
		out     gradle.UpdatePreconditionOut
	}{
		{
			name:    "up-to-date",
			content: testdata.GradleWrapperProperties50,
			out:     gradle.AlreadyHasLatestGradle,
		},
		{
			name:    "out-of-date",
			content: testdata.GradleWrapperProperties4102,
			out:     gradle.ReadyToUpdate,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
			repositoryRepository.EXPECT().Get(ctx, repositoryID).
				Return(&git.Repository{}, nil)
			repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, gradle.WrapperPropertiesPath).
				Return(c.content, nil)
			repositoryRepository.EXPECT().GetReadme(ctx, repositoryID).
				Return(readme, nil)

			gradleService := gateways.NewMockGradleService(ctrl)
			gradleService.EXPECT().GetCurrentRelease(ctx).
				Return(&gradle.Release{Version: "5.0"}, nil)

			u := GetRepository{
				RepositoryRepository: repositoryRepository,
				GradleService:        gradleService,
			}
			resp, err := u.Do(ctx, repositoryID)
			if err != nil {
				t.Fatalf("could not do usecase: %s", err)
			}
			if resp.GradleUpdatePreconditionOut != c.out {
				t.Errorf("UpdatePreconditionOut wants %v but %v", c.out, resp.GradleUpdatePreconditionOut)
			}
		})
	}
}

func TestGetRepository_Do_NoSuchRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	repositoryID := git.RepositoryID{Owner: "owner", Name: "repo"}

	repositoryError := gateways.NewMockRepositoryError(ctrl)
	repositoryError.EXPECT().NoSuchEntity().AnyTimes().Return(true)

	repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
	repositoryRepository.EXPECT().Get(ctx, repositoryID).
		Return(nil, repositoryError)

	gradleService := gateways.NewMockGradleService(ctrl)

	u := GetRepository{
		RepositoryRepository: repositoryRepository,
		GradleService:        gradleService,
	}
	resp, err := u.Do(ctx, repositoryID)
	if resp != nil {
		t.Errorf("resp wants nil but %+v", resp)
	}
	if err == nil {
		t.Fatalf("err wants non-nil but nil")
	}
	cause, ok := errors.Cause(err).(usecases.GetRepositoryError)
	if !ok {
		t.Fatalf("cause wants GetRepositoryError but %T", cause)
	}
	if cause.NoSuchRepository() != true {
		t.Errorf("NoSuchRepository wants true")
	}
}
