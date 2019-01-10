package usecases_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/testdata"
	"github.com/int128/gradleupdate/gateways/interfaces/mock_gateways"
	"github.com/int128/gradleupdate/usecases"
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
			content:        domain.FileContent(testdata.GradleWrapperProperties4102),
			currentVersion: "4.10.2",
			latestVersion:  "4.10.2",
			upToDate:       true,
		},
		{
			name:           "out-of-date",
			content:        domain.FileContent(testdata.GradleWrapperProperties4102),
			currentVersion: "4.10.2",
			latestVersion:  "5.1",
			upToDate:       false,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			repositoryRepository := mock_gateways.NewMockRepositoryRepository(ctrl)
			repositoryRepository.EXPECT().Get(ctx, repositoryID).
				Return(&domain.Repository{}, nil)
			repositoryRepository.EXPECT().GetFileContent(ctx, repositoryID, domain.GradleWrapperPropertiesPath).
				Return(c.content, nil)

			gradleService := mock_gateways.NewMockGradleService(ctrl)
			gradleService.EXPECT().GetCurrentVersion(ctx).
				Return(c.latestVersion, nil)

			u := usecases.GetRepository{
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
