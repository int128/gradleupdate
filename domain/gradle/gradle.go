package gradle

import (
	"regexp"

	"github.com/int128/gradleupdate/domain/git"
)

type Release struct {
	Version Version
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
