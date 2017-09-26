package org.hidetake.gradleupdate.domain

import org.eclipse.egit.github.core.User

interface UserRepository {
    fun getLoginUser(): User
}
