package org.hidetake.gradleupdate.domain

object GradleWrapperPullRequestFactory {
    fun create(
        repositoryName: String,
        gradleWrapperVersion: GradleWrapperVersion,
        files: List<GradleWrapperFile>
    ) =
        GradleWrapperPullRequest(
            "Gradle ${gradleWrapperVersion.version}",
            "Gradle ${gradleWrapperVersion.version} is available.",
            repositoryName,
            branchName(repositoryName, gradleWrapperVersion),
            "Gradle Update",
            "gradleupdate@users.noreply.github.com"
            , files)

    fun branchName(repositoryName: String, gradleWrapperVersion: GradleWrapperVersion) =
        "gradle-${gradleWrapperVersion.version}-${repositoryName.split("/").first()}"
}
