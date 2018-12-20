package usecases

import "github.com/int128/gradleupdate/domain"

const gradleWrapperPropertiesPath = "gradle/wrapper/gradle-wrapper.properties"

var gradleWrapperFiles = []domain.File{
	{
		Path: gradleWrapperPropertiesPath,
		Mode: "100644",
	},
	{
		Path: "gradle/wrapper/gradle-wrapper.jar",
		Mode: "100644",
	},
	{
		Path: "gradlew",
		Mode: "100755",
	},
	{
		Path: "gradlew.bat",
		Mode: "100644",
	},
}

func findGradleWrapperVersion(files []domain.File) domain.GradleVersion {
	for _, file := range files {
		if file.Path == gradleWrapperPropertiesPath {
			return domain.FindGradleWrapperVersion(string(file.Content))
		}
	}
	return ""
}
