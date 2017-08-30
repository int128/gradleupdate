package org.hidetake.gradleupdate.domain

object GradleWrapperPullRequestFactory {
    fun create(gradleWrapperVersion: GradleWrapperVersion, files: List<GradleWrapperFile>): GradleWrapperPullRequest =
        GradleWrapperPullRequest(
            "Gradle ${gradleWrapperVersion.version}",
            "Gradle ${gradleWrapperVersion.version}",
            "gradle-${gradleWrapperVersion.version}",
            files
        )
}
