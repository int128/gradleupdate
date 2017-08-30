package org.hidetake.gradleupdate.domain

interface PullRequestRepository {
    fun create(repositoryName: String, gradleWrapperPullRequest: GradleWrapperPullRequest)
}
