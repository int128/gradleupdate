package service

import (
	"context"
	"fmt"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/infrastructure"
)

// GradleWrapperStatus represents whether the wrapper is up-to-date or out-of-date.
type GradleWrapperStatus struct {
	TargetVersion domain.GradleVersion
	LatestVersion domain.GradleVersion
	UpToDate      bool
}

// GetGradleWrapperStatus returns a GradleWrapperStatus of the repository.
func GetGradleWrapperStatus(ctx context.Context, owner, repo string) (*GradleWrapperStatus, error) {
	targetVersion, err := GetGradleWrapperVersion(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("Could not get version of %s/%s: %s", owner, repo, err)
	}
	latestVersion, err := GetGradleWrapperVersion(ctx, "int128", "latest-gradle-wrapper")
	if err != nil {
		return nil, fmt.Errorf("Could not get the latest version: %s", err)
	}
	return &GradleWrapperStatus{
		TargetVersion: targetVersion,
		LatestVersion: latestVersion,
		UpToDate:      domain.IsUpToDate(targetVersion, latestVersion),
	}, nil
}

const gradleWrapperPropsPath = "gradle/wrapper/gradle-wrapper.properties"

// GetGradleWrapperVersion returns version of the wrapper in the repository.
func GetGradleWrapperVersion(ctx context.Context, owner, repo string) (domain.GradleVersion, error) {
	c := infrastructure.GitHubClient(ctx)
	fc, _, _, err := c.Repositories.GetContents(ctx, owner, repo, gradleWrapperPropsPath, nil)
	if err != nil {
		return "", fmt.Errorf("Could not get content: %s", err)
	}
	if fc == nil {
		return "", fmt.Errorf("No such file: %s", gradleWrapperPropsPath)
	}
	content, err := fc.GetContent()
	if err != nil {
		return "", fmt.Errorf("Could not decode content: %s", err)
	}
	v := domain.FindGradleWrapperVersion(content)
	if v == "" {
		return "", fmt.Errorf("Could not determine version from gradle-wrapper.properties")
	}
	return v, nil
}
