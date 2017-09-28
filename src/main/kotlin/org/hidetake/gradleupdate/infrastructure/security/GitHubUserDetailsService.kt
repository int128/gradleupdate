package org.hidetake.gradleupdate.infrastructure.security

import com.google.appengine.api.memcache.MemcacheService
import org.springframework.security.core.userdetails.AuthenticationUserDetailsService
import org.springframework.security.core.userdetails.User
import org.springframework.security.core.userdetails.UserDetails
import org.springframework.security.web.authentication.preauth.PreAuthenticatedAuthenticationToken

class GitHubUserDetailsService(private val memcacheService: MemcacheService)
    : AuthenticationUserDetailsService<PreAuthenticatedAuthenticationToken> {
    override fun loadUserDetails(token: PreAuthenticatedAuthenticationToken): UserDetails {
        val sessionToken = token.principal as? String
        return User(sessionToken, "", arrayListOf(GradleUpdateUser()))
    }
}
