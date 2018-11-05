package domain

import "regexp"

var regexpGradleWrapperVersion = regexp.MustCompile("distributionUrl=.+?/gradle-(.+?)-.+?\\.zip")

func FindGradleWrapperVersion(gradleWrapperProperties string) GradleVersion {
	m := regexpGradleWrapperVersion.FindStringSubmatch(gradleWrapperProperties)
	if len(m) != 2 {
		return ""
	}
	return GradleVersion(m[1])
}
