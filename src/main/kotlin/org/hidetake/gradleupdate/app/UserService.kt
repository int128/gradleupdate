package org.hidetake.gradleupdate.app

import org.hidetake.gradleupdate.domain.UserRepository
import org.springframework.stereotype.Service

@Service
class UserService(private val userRepository: UserRepository) {
    fun getLoginUser() = userRepository.getLoginUser()
}
