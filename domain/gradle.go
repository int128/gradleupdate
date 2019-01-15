package domain

import (
	"regexp"
)

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
