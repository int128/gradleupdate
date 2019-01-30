package usecases

import (
	"context"
	"fmt"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/gateways/interfaces"
	"github.com/int128/gradleupdate/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// SendUpdate provides a use case to send a pull request for updating Gradle in a repository.
type SendUpdate struct {
	dig.In
	GradleService                gateways.GradleService
	RepositoryRepository         gateways.RepositoryRepository
	RepositoryLastScanRepository gateways.RepositoryLastScanRepository
	SendPullRequest              usecases.SendPullRequest
	TimeService                  gateways.TimeService
}

func (usecase *SendUpdate) Do(ctx context.Context, id domain.RepositoryID, badgeURL string) error {
	scan := domain.RepositoryLastScan{
		Repository:   id,
		LastScanTime: usecase.TimeService.Now(),
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
	//TODO: fetch entities concurrently
	readme, err := usecase.RepositoryRepository.GetReadme(ctx, id)
	if err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			if err.NoSuchEntity() {
				return errors.Wrapf(&sendUpdateError{error: err, noReadmeBadge: true}, "README did not found")
			}
		}
		return errors.Wrapf(err, "error while getting README")
	}
	gradleWrapperProperties, err := usecase.RepositoryRepository.GetFileContent(ctx, id, domain.GradleWrapperPropertiesPath)
	if err != nil {
		if err, ok := errors.Cause(err).(gateways.RepositoryError); ok {
			if err.NoSuchEntity() {
				return errors.Wrapf(&sendUpdateError{error: err, noGradleVersion: true}, "gradle-wrapper.properties did not found")
			}
		}
		return errors.Wrapf(err, "error while getting gradle-wrapper.properties")
	}
	latestRelease, err := usecase.GradleService.GetCurrentRelease(ctx)
	if err != nil {
		return errors.Wrapf(err, "error while getting the latest Gradle version")
	}

	out := domain.CheckGradleUpdatePrecondition(domain.GradleUpdatePreconditionIn{
		Readme:                  readme,
		BadgeURL:                badgeURL,
		GradleWrapperProperties: gradleWrapperProperties,
		LatestGradleRelease:     latestRelease,
	})
	if out.NoReadmeBadge {
		return errors.WithStack(&sendUpdateError{error: fmt.Errorf("README did not contain the badge"), noReadmeBadge: true})
	}
	if out.NoGradleVersion {
		return errors.WithStack(&sendUpdateError{error: fmt.Errorf("properties did not contain version string"), noGradleVersion: true})
	}
	if out.AlreadyHasLatestGradle {
		return errors.WithStack(&sendUpdateError{error: fmt.Errorf("current version is already latest"), alreadyHasLatestGradle: true})
	}

	newProps := domain.ReplaceGradleWrapperVersion(gradleWrapperProperties, latestRelease.Version)
	req := usecases.SendPullRequestRequest{
		Base:           id,
		HeadBranchName: fmt.Sprintf("gradle-%s-%s", latestRelease.Version, id.Owner),
		CommitMessage:  fmt.Sprintf("Gradle %s", latestRelease.Version),
		CommitFiles: []domain.File{
			{
				Path:    domain.GradleWrapperPropertiesPath,
				Content: domain.FileContent(newProps),
			},
		},
		Title: fmt.Sprintf("Gradle %s", latestRelease.Version),
		Body:  fmt.Sprintf(`Gradle %s is available.`, latestRelease.Version),
	}
	if err := usecase.SendPullRequest.Do(ctx, req); err != nil {
		return errors.Wrapf(err, "error while sending a pull request %+v", req)
	}
	return nil
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
