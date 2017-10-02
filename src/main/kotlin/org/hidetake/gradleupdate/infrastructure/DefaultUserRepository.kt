package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.service.UserService
import org.hidetake.gradleupdate.domain.UserRepository
import org.springframework.stereotype.Repository

@Repository
class DefaultUserRepository(loginUserGitHubClient: LoginUserGitHubClient) : UserRepository {
    private val userService = UserService(loginUserGitHubClient)

    override fun getLoginUser() = userService.user!!
}
