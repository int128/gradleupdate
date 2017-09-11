package org.hidetake.gradleupdate.domain

interface PullRequestRepository {
    fun createOrUpdate(
        repositoryPath: RepositoryPath,
        gradleWrapperVersion: GradleWrapperVersion,
        files: List<GradleWrapperFile>
    )

    fun find(
        repositoryPath: RepositoryPath,
        gradleWrapperVersion: GradleWrapperVersion
    ): PullRequestForUpdate?
}
