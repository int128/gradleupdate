package usecases

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
)

// SendUpdate provides a use case to send a pull request for updating Gradle in a repository.
type SendUpdate struct {
	GradleService                gateways.GradleService
	RepositoryRepository         gateways.RepositoryRepository
	RepositoryLastScanRepository gateways.RepositoryLastScanRepository
	SendPullRequest              usecases.SendPullRequest
	NowFunc                      func() time.Time
}

func (usecase *SendUpdate) Do(ctx context.Context, id domain.RepositoryID, badgeURL string) error {
	scan := domain.RepositoryLastScan{
		Repository:   id,
		LastScanTime: usecase.Now(),
	}
	err := usecase.sendUpdate(ctx, id, badgeURL)
	if err != nil {
		if err, ok := errors.Cause(err).(*sendUpdateError); ok {
			scan.NoGradleVersionError = err.noGradleVersion
			scan.NoReadmeBadgeError = err.noReadmeBadge
			scan.AlreadyLatestGradleError = err.alreadyHasLatestGradle
		}
	}
	if err := usecase.RepositoryLastScanRepository.Save(ctx, scan); err != nil {
		return errors.Wrapf(err, "error while saving the scan for the repository %s", id)
	}
	return errors.Wrapf(err, "error while scanning the repository %s", id)
}

func (usecase *SendUpdate) sendUpdate(ctx context.Context, id domain.RepositoryID, badgeURL string) error {
	latestVersion, err := usecase.GradleService.GetCurrentVersion(ctx)
	if err != nil {
		return errors.Wrapf(err, "error while getting the latest Gradle version")
	}

	// check if the properties file has out-of-date version
	gradleWrapperProperties, err := usecase.RepositoryRepository.GetFileContent(ctx, id, domain.GradleWrapperPropertiesPath)
	if err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			if err.NoSuchEntity() {
				return errors.Wrapf(&sendUpdateError{error: err, noGradleVersion: true}, "could not find Gradle version")
			}
		}
		return errors.Wrapf(err, "error while getting the properties file")
	}
	currentVersion := domain.FindGradleWrapperVersion(gradleWrapperProperties)
	if currentVersion == "" {
		return errors.WithStack(&sendUpdateError{error: fmt.Errorf("properties did not contain version string"), noGradleVersion: true})
	}
	if currentVersion.GreaterOrEqualThan(latestVersion) {
		return errors.WithStack(&sendUpdateError{error: fmt.Errorf("current version %s is already latest", currentVersion), alreadyHasLatestGradle: true})
	}

	// check if the README has the badge
	readme, err := usecase.RepositoryRepository.GetReadme(ctx, id)
	if err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			if err.NoSuchEntity() {
				return errors.Wrapf(&sendUpdateError{error: err, noReadmeBadge: true}, "could not find README")
			}
		}
		return errors.Wrapf(err, "error while getting README")
	}
	if !bytes.Contains(readme, []byte(badgeURL)) {
		return errors.WithStack(&sendUpdateError{error: fmt.Errorf("README did not contain the badge"), noReadmeBadge: true})
	}

	// send a pull request for the latest version
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
	if err := usecase.SendPullRequest.Do(ctx, req); err != nil {
		return errors.Wrapf(err, "error while sending a pull request %+v", req)
	}
	return nil
}

func (usecase *SendUpdate) Now() time.Time {
	if usecase.NowFunc != nil {
		return usecase.NowFunc()
	}
	return time.Now()
}

type sendUpdateError struct {
	error
	noGradleVersion        bool
	noReadmeBadge          bool
	alreadyHasLatestGradle bool
}

func (err *sendUpdateError) NoGradleVersion() bool        { return err.noGradleVersion }
func (err *sendUpdateError) NoReadmeBadge() bool          { return err.noReadmeBadge }
func (err *sendUpdateError) AlreadyHasLatestGradle() bool { return err.alreadyHasLatestGradle }
