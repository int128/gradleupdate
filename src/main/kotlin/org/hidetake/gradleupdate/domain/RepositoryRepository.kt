package org.hidetake.gradleupdate.domain

import org.eclipse.egit.github.core.Repository

interface RepositoryRepository {
    fun getByName(repositoryPath: RepositoryPath): Repository
}
