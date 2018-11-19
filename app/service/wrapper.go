package service

import (
	"context"
	"encoding/base64"

	"github.com/google/go-github/v18/github"
	"github.com/int128/gradleupdate/app/service/pr"
	"github.com/pkg/errors"
)

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
