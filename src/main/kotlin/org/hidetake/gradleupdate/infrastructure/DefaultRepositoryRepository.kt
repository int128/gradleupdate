package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.Repository
import org.eclipse.egit.github.core.client.GitHubClient
import org.eclipse.egit.github.core.service.RepositoryService
import org.hidetake.gradleupdate.domain.RepositoryRepository

@org.springframework.stereotype.Repository
class DefaultRepositoryRepository(private val client: GitHubClient) : RepositoryRepository {
    override fun getByName(name: String): Repository =
        RepositoryService(client).getRepository({name})
}
