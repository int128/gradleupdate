package org.hidetake.gradleupdate.domain

import org.eclipse.egit.github.core.PullRequest

interface PullRequestRepository {
    fun createOrUpdate(
        repositoryPath: RepositoryPath,
        gradleWrapperVersion: GradleWrapperVersion,
        files: List<GradleWrapperFile>
    )

    fun find(
        repositoryPath: RepositoryPath,
        gradleWrapperVersion: GradleWrapperVersion
    ): PullRequest?
}
