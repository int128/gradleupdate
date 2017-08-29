package org.hidetake.gradleupdate.domain

interface GradleWrapperRepository {
    fun findVersion(repositoryName: String): GradleWrapperVersion?
}
