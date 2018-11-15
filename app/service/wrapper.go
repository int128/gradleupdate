package service

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/app/domain"
	"github.com/int128/gradleupdate/app/infrastructure"
	"github.com/int128/gradleupdate/app/service/pr"
	"github.com/pkg/errors"
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

var gradleWrapperFiles = []pr.File{
	pr.File{
		Path: "gradle/wrapper/gradle-wrapper.properties",
		Mode: "100644",
	},
	pr.File{
		Path: "gradle/wrapper/gradle-wrapper.jar",
		Mode: "100644",
	},
	pr.File{
		Path: "gradlew",
		Mode: "100755",
	},
	pr.File{
		Path: "gradlew.bat",
		Mode: "100644",
	},
}

// FindGradleWrapperFiles returns files of the latest Gradle wrapper.
func FindGradleWrapperFiles(ctx context.Context, c *github.Client, owner, repo string) ([]pr.File, error) {
	r := make([]pr.File, len(gradleWrapperFiles))
	for i, file := range gradleWrapperFiles {
		fc, _, _, err := c.Repositories.GetContents(ctx, owner, repo, file.Path, &github.RepositoryContentGetOptions{})
		if err != nil {
			return nil, errors.Wrapf(err, "Could not get content of file %s", file.Path)
		}
		content, err := fc.GetContent()
		if err != nil {
			return nil, errors.Wrapf(err, "Could not decode content of file %s", file.Path)
		}
		r[i] = file
		r[i].EncodedContent = base64.StdEncoding.EncodeToString([]byte(content))
	}
	return r, nil
}
