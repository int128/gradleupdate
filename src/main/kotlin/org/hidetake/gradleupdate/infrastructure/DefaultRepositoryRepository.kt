package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.Repository
import org.eclipse.egit.github.core.service.RepositoryService
import org.hidetake.gradleupdate.domain.RepositoryPath
import org.hidetake.gradleupdate.domain.RepositoryRepository

@org.springframework.stereotype.Repository
class DefaultRepositoryRepository(
    systemGitHubClient: SystemGitHubClient,
    loginUserGitHubClient: LoginUserGitHubClient
) : RepositoryRepository {
    private val systemRepositoryService = RepositoryService(systemGitHubClient)
    private val loginUserRepositoryService = RepositoryService(loginUserGitHubClient)

    override fun getByName(repositoryPath: RepositoryPath): Repository =
        systemRepositoryService.getRepository({repositoryPath.fullName})

    override fun findAllOfLoginUser(): List<Repository> =
        loginUserRepositoryService.repositories
}
