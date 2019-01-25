package usecases

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/testdata"
	"github.com/int128/gradleupdate/gateways/interfaces/test_doubles"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
)

func TestGetRepository_Do(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	repositoryID := domain.RepositoryID{Owner: "owner", Name: "repo"}

	for _, c := range []struct {
		name           string
		content        domain.FileContent
		currentVersion domain.GradleVersion
		latestVersion  domain.GradleVersion
		upToDate       bool
	}{
		{
			name:           "up-to-date",
			content:        testdata.GradleWrapperProperties4102,
			currentVersion: "4.10.2",
			latestVersion:  "4.10.2",
			upToDate:       true,
		},
		{
			name:           "out-of-date",
			content:        testdata.GradleWrapperProperties4102,
			currentVersion: "4.10.2",
			latestVersion:  "5.1",
			upToDate:       false,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			repositoryRepository := gateways.NewMockRepositoryRepository(ctrl)
			repositoryRepository.EXPECT().Get(ctx, repositoryID).
				Return(&domain.Repository{}, nil)
			repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
				Return(c.content, nil)

			gradleService := gateways.NewMockGradleService(ctrl)
			gradleService.EXPECT().GetCurrentVersion(ctx).
				Return(c.latestVersion, nil)

			u := GetRepository{
				RepositoryRepository: repositoryRepository,
				GradleService:        gradleService,
			}
			resp, err := u.Do(ctx, repositoryID)
			if err != nil {
				t.Fatalf("could not do usecase: %s", err)
			}
			if resp.CurrentVersion != c.currentVersion {
				t.Errorf("CurrentVersion wants %s but %s", c.currentVersion, resp.CurrentVersion)
			}
			if resp.LatestVersion != c.latestVersion {
				t.Errorf("LatestVersion wants %s but %s", c.latestVersion, resp.LatestVersion)
			}
			if resp.UpToDate != c.upToDate {
				t.Errorf("UpToDate wants %v but %v", c.upToDate, resp.UpToDate)
			}
		})
	}
}

func TestGetRepository_Do_NoSuchRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	repositoryID := domain.RepositoryID{Owner: "owner", Name: "repo"}

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
	if cause.NoGradleVersion() != false {
		t.Errorf("NoGradleVersion wants false")
	}
	if cause.NoSuchRepository() != true {
		t.Errorf("NoSuchRepository wants true")
	}
}
