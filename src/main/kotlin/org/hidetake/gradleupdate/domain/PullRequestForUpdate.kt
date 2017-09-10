package org.hidetake.gradleupdate.domain

class PullRequestForUpdate(
    val title: String,
    val description: String,
    val repositoryName: String,
    val branch: Branch,
    val files: List<GradleWrapperFile> = emptyList()
)
