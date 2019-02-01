package gradle

import (
	"bytes"
	"regexp"

	"github.com/int128/gradleupdate/domain/git"
)

type Release struct {
	Version Version
}

type UpdatePreconditionIn struct {
	Readme                  git.FileContent
	BadgeURL                string
	GradleWrapperProperties git.FileContent
	LatestGradleRelease     *Release
}

type UpdatePreconditionOut int

const (
	ReadyToUpdate             = UpdatePreconditionOut(0)
	AlreadyHasLatestGradle    = UpdatePreconditionOut(1)
	NoGradleWrapperProperties = UpdatePreconditionOut(51)
	NoGradleVersion           = UpdatePreconditionOut(52)
	NoReadme                  = UpdatePreconditionOut(53)
	NoReadmeBadge             = UpdatePreconditionOut(54)
)

func CheckUpdatePrecondition(in UpdatePreconditionIn) UpdatePreconditionOut {
	if in.GradleWrapperProperties == nil {
		return NoGradleWrapperProperties
	}
	currentGradleVersion := FindWrapperVersion(in.GradleWrapperProperties)
	if currentGradleVersion == "" {
		return NoGradleVersion
	}
	if in.Readme == nil {
		return NoReadme
	}
	if !bytes.Contains(in.Readme, []byte(in.BadgeURL)) {
		return NoReadmeBadge
	}
	if currentGradleVersion.GreaterOrEqualThan(in.LatestGradleRelease.Version) {
		return AlreadyHasLatestGradle
	}
	return ReadyToUpdate
}

// WrapperPropertiesPath is path to the gradle-wrapper.properties
const WrapperPropertiesPath = "gradle/wrapper/gradle-wrapper.properties"

var regexpWrapperVersion = regexp.MustCompile(`(distributionUrl=.+?/gradle-)(.+?)(-.+?\.zip)`)

// FindWrapperVersion returns Gradle version in a properties file.
// It returns an empty string if version does not find.
func FindWrapperVersion(gradleWrapperProperties git.FileContent) Version {
	m := regexpWrapperVersion.FindSubmatch(gradleWrapperProperties)
	if len(m) != 4 {
		return ""
	}
	return Version(m[2])
}

// ReplaceWrapperVersion returns content with the given version.
func ReplaceWrapperVersion(gradleWrapperProperties git.FileContent, version Version) git.FileContent {
	return regexpWrapperVersion.ReplaceAll(gradleWrapperProperties,
		[]byte("${1}"+version.String()+"${3}"))
}
