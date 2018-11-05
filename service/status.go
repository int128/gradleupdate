package service

import (
	"context"
	"fmt"

	"github.com/int128/gradleupdate/domain"
)

type GradleWrapperStatus struct {
	TargetVersion domain.GradleVersion
	LatestVersion domain.GradleVersion
	UpToDate      bool
}

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
