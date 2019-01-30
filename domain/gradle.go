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

type GradleUpdatePreconditionOut struct {
	NoReadmeBadge          bool
	NoGradleVersion        bool
	AlreadyHasLatestGradle bool
}

func CheckGradleUpdatePrecondition(in GradleUpdatePreconditionIn) (out GradleUpdatePreconditionOut) {
	if !bytes.Contains(in.Readme, []byte(in.BadgeURL)) {
		out.NoReadmeBadge = true
	}
	currentVersion := FindGradleWrapperVersion(in.GradleWrapperProperties)
	if currentVersion == "" {
		out.NoGradleVersion = true
	} else {
		if currentVersion.GreaterOrEqualThan(in.LatestGradleRelease.Version) {
			out.AlreadyHasLatestGradle = true
		}
	}
	return
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
