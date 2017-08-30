package org.hidetake.gradleupdate.domain

class GradleWrapperPullRequest(
    val title: String,
    val description: String,
    val branchName: String,
    val files: List<GradleWrapperFile>
)
