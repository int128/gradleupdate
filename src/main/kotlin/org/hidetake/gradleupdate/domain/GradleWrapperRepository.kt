package org.hidetake.gradleupdate.domain

interface GradleWrapperRepository {
    fun findVersion(repositoryPath: RepositoryPath): GradleWrapperVersion?

    fun findFiles(repositoryPath: RepositoryPath): List<GradleWrapperFile>
}
