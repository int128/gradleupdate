package org.hidetake.gradleupdate.domain

interface PullRequestRepository {
    fun createOrUpdate(gradleWrapperPullRequest: GradleWrapperPullRequest): Unit

    fun find(repositoryName: String, version: GradleWrapperVersion): GradleWrapperPullRequest?
}
