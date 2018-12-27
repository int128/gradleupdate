package domain

import "regexp"

// GradleWrapperPropertiesPath is path to the gradle-wrapper.properties
const GradleWrapperPropertiesPath = "gradle/wrapper/gradle-wrapper.properties"

var regexpGradleWrapperVersion = regexp.MustCompile("distributionUrl=.+?/gradle-(.+?)-.+?\\.zip")

// FindGradleWrapperVersion returns Gradle version in a properties file.
func FindGradleWrapperVersion(gradleWrapperProperties string) GradleVersion {
	m := regexpGradleWrapperVersion.FindStringSubmatch(gradleWrapperProperties)
	if len(m) != 2 {
		return ""
	}
	return GradleVersion(m[1])
}
