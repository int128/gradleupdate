package org.hidetake.gradleupdate.domain

import org.eclipse.egit.github.core.PullRequest

interface PullRequestRepository {
    fun createOrUpdate(gradleWrapperPullRequest: GradleWrapperPullRequest): Unit

    fun find(repositoryName: String, version: GradleWrapperVersion): PullRequest?
}
