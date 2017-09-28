package org.hidetake.gradleupdate.infrastructure.security

import org.springframework.security.core.GrantedAuthority

class GradleUpdateUser : GrantedAuthority {
    override fun getAuthority(): String = "ROLE_USER"
}
