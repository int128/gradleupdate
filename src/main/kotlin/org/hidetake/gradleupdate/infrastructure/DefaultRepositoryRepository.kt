package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.Repository
import org.eclipse.egit.github.core.service.RepositoryService
import org.hidetake.gradleupdate.domain.RepositoryPath
import org.hidetake.gradleupdate.domain.RepositoryRepository
import org.hidetake.gradleupdate.infrastructure.egit.EnhancedGitHubClient
import org.springframework.beans.factory.annotation.Qualifier

@org.springframework.stereotype.Repository
class DefaultRepositoryRepository(@Qualifier("gradleUpdateGitHubClient") client: EnhancedGitHubClient) : RepositoryRepository {
    private val repositoryService = RepositoryService(client)

    override fun getByName(repositoryPath: RepositoryPath): Repository =
        repositoryService.getRepository({repositoryPath.fullName})
}
