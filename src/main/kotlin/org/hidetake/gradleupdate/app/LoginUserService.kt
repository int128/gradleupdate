package org.hidetake.gradleupdate.app

import org.hidetake.gradleupdate.domain.RepositoryRepository
import org.hidetake.gradleupdate.domain.UserRepository
import org.springframework.stereotype.Service

@Service
class LoginUserService(
    private val userRepository: UserRepository,
    private val repositoryRepository: RepositoryRepository
) {
    fun getLoginUser() = userRepository.getLoginUser()

    fun getRepositories() = repositoryRepository.findAllOfLoginUser()
}
