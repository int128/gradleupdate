package domain

import (
	"bytes"
	"regexp"
)

type GradleRelease struct {
	Version GradleVersion
}

type GradleUpdatePreconditionIn struct {
	Readme                  FileContent
	BadgeURL                string
	GradleWrapperProperties FileContent
	LatestGradleRelease     *GradleRelease
}

type GradleUpdatePreconditionOut int

var (
	ReadyToUpdate             = GradleUpdatePreconditionOut(0)
	AlreadyHasLatestGradle    = GradleUpdatePreconditionOut(1)
	NoGradleWrapperProperties = GradleUpdatePreconditionOut(51)
	NoGradleVersion           = GradleUpdatePreconditionOut(52)
	NoReadme                  = GradleUpdatePreconditionOut(53)
	NoReadmeBadge             = GradleUpdatePreconditionOut(54)
)

func CheckGradleUpdatePrecondition(in GradleUpdatePreconditionIn) GradleUpdatePreconditionOut {
	if in.GradleWrapperProperties == nil {
		return NoGradleWrapperProperties
	}
	currentGradleVersion := FindGradleWrapperVersion(in.GradleWrapperProperties)
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

// GradleWrapperPropertiesPath is path to the gradle-wrapper.properties
const GradleWrapperPropertiesPath = "gradle/wrapper/gradle-wrapper.properties"

var regexpGradleWrapperVersion = regexp.MustCompile(`(distributionUrl=.+?/gradle-)(.+?)(-.+?\.zip)`)

// FindGradleWrapperVersion returns Gradle version in a properties file.
// It returns an empty string if version does not find.
func FindGradleWrapperVersion(gradleWrapperProperties FileContent) GradleVersion {
	m := regexpGradleWrapperVersion.FindSubmatch(gradleWrapperProperties)
	if len(m) != 4 {
		return ""
	}
	return GradleVersion(m[2])
}

// ReplaceGradleWrapperVersion returns content with the given version.
func ReplaceGradleWrapperVersion(gradleWrapperProperties FileContent, version GradleVersion) FileContent {
	return regexpGradleWrapperVersion.ReplaceAll(gradleWrapperProperties,
		[]byte("${1}"+version.String()+"${3}"))
}
