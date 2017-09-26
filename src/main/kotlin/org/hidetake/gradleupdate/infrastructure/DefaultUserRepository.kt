package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.service.UserService
import org.hidetake.gradleupdate.domain.UserRepository
import org.hidetake.gradleupdate.infrastructure.egit.EnhancedGitHubClient
import org.springframework.beans.factory.annotation.Qualifier
import org.springframework.stereotype.Repository

@Repository
class DefaultUserRepository(
    @Qualifier("contextGitHubClient") contextGitHubClient: EnhancedGitHubClient
) : UserRepository {
    private val userService = UserService(contextGitHubClient)

    override fun getLoginUser() = userService.user!!
}
