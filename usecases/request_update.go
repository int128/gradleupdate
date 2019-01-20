package usecases

import (
	"bytes"
	"context"
	"fmt"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
)

// RequestUpdate provides a use case to send a pull request for updating Gradle in a repository.
type RequestUpdate struct {
	GradleService        gateways.GradleService
	RepositoryRepository gateways.RepositoryRepository
	SendPullRequest      usecases.SendPullRequest
}

func (usecase *RequestUpdate) Do(ctx context.Context, id domain.RepositoryID, badgeURL string) error {
	latestVersion, err := usecase.GradleService.GetCurrentVersion(ctx)
	if err != nil {
		return errors.Wrapf(err, "could not get the latest Gradle version")
	}
	readme, err := usecase.RepositoryRepository.GetReadme(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "could not find readme file")
	}
	if !bytes.Contains(readme, []byte(badgeURL)) {
		return errors.Wrapf(&requestUpdateError{noBadgeInReadme: true}, "readme did not contain any badge URL")
	}
	gradleWrapperProperties, err := usecase.RepositoryRepository.GetFileContent(ctx, id, domain.GradleWrapperPropertiesPath)
	if err != nil {
		return errors.Wrapf(err, "could not find properties file")
	}
	currentVersion := domain.FindGradleWrapperVersion(gradleWrapperProperties)
	if currentVersion == "" {
		return errors.Errorf("could not find version in the properties")
	}
	if currentVersion == latestVersion {
		return nil // branch is already up-to-date
	}
	newProps := domain.ReplaceGradleWrapperVersion(gradleWrapperProperties, latestVersion)

	req := usecases.SendPullRequestRequest{
		Base:           id,
		HeadBranchName: fmt.Sprintf("gradle-%s-%s", latestVersion, id.Owner),
		CommitMessage:  fmt.Sprintf("Gradle %s", latestVersion),
		CommitFiles: []domain.File{
			{
				Path:    domain.GradleWrapperPropertiesPath,
				Content: domain.FileContent(newProps),
			},
		},
		Title: fmt.Sprintf("Gradle %s", latestVersion),
		Body:  fmt.Sprintf(`Gradle %s is available.`, latestVersion),
	}
	if err = usecase.SendPullRequest.Do(ctx, req); err != nil {
		return errors.Wrapf(err, "could not send a pull request %+v", req)
	}
	return nil
}

type requestUpdateError struct {
	error
	noBadgeInReadme bool
}

func (err *requestUpdateError) NoBadgeInReadme() bool { return err.noBadgeInReadme }
