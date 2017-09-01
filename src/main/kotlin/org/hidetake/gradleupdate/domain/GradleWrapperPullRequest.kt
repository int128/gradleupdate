package org.hidetake.gradleupdate.domain

class GradleWrapperPullRequest(
    val title: String,
    val description: String,
    val repositoryName: String,
    val branchName: String,
    val authorName: String,
    val authorEmail: String,
    val files: List<GradleWrapperFile>
)
