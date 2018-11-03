package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/go-github/v18/github"
)

const gradleWrapperPropsPath = "gradle/wrapper/gradle-wrapper.properties"

func GetGradleWrapperVersion(ctx context.Context, owner, repo string) (string, error) {
	c := github.NewClient(nil)

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
	v := findGradleWrapperVersion(content)
	if v == "" {
		return "", fmt.Errorf("Could not determine version from file")
	}
	return v, nil
}

var regexpGradleWrapperVersion = regexp.MustCompile("distributionUrl=.+?/gradle-(.+?)-.+?\\.zip")

func findGradleWrapperVersion(content string) string {
	m := regexpGradleWrapperVersion.FindStringSubmatch(content)
	if len(m) != 2 {
		return ""
	}
	return m[1]
}
