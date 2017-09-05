package org.hidetake.gradleupdate.domain

class GradleWrapperPullRequest(
    val title: String,
    val description: String,
    val repositoryName: String,
    val branch: GradleWrapperPullRequestBranch,
    val files: List<GradleWrapperFile> = emptyList()
) {
    object Factory {
        fun create(
            repositoryName: String,
            gradleWrapperVersion: GradleWrapperVersion,
            files: List<GradleWrapperFile>
        ) =
            GradleWrapperPullRequest(
                "Gradle ${gradleWrapperVersion.version}",
                "Gradle ${gradleWrapperVersion.version} is available.",
                repositoryName,
                GradleWrapperPullRequestBranch.Factory.create(repositoryName, gradleWrapperVersion),
                files)
    }
}
