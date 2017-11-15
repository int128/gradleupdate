package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.Repository
import org.eclipse.egit.github.core.service.RepositoryService
import org.hidetake.gradleupdate.domain.RepositoryPath
import org.hidetake.gradleupdate.domain.RepositoryRepository

@org.springframework.stereotype.Repository
class DefaultRepositoryRepository(client: SystemGitHubClient) : RepositoryRepository {
    private val repositoryService = RepositoryService(client)

    override fun getByName(repositoryPath: RepositoryPath): Repository =
        repositoryService.getRepository({repositoryPath.fullName})
}
