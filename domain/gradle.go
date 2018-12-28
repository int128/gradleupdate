package domain

import (
	"regexp"
)

// GradleWrapperPropertiesPath is path to the gradle-wrapper.properties
const GradleWrapperPropertiesPath = "gradle/wrapper/gradle-wrapper.properties"

var regexpGradleWrapperVersion = regexp.MustCompile(`(distributionUrl=.+?/gradle-)(.+?)(-.+?\.zip)`)

// FindGradleWrapperVersion returns Gradle version in a properties file.
func FindGradleWrapperVersion(gradleWrapperProperties string) GradleVersion {
	m := regexpGradleWrapperVersion.FindStringSubmatch(gradleWrapperProperties)
	if len(m) != 4 {
		return ""
	}
	return GradleVersion(m[2])
}

// ReplaceGradleWrapperVersion returns content with the given version.
func ReplaceGradleWrapperVersion(gradleWrapperProperties string, version GradleVersion) string {
	return regexpGradleWrapperVersion.ReplaceAllString(gradleWrapperProperties,
		"${1}"+version.String()+"${3}")
}
