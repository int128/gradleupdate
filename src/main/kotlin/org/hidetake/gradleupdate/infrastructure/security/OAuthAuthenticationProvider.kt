package org.hidetake.gradleupdate.infrastructure.security

import org.springframework.security.authentication.AuthenticationProvider
import org.springframework.security.core.Authentication
import org.springframework.stereotype.Component

@Component
class OAuthAuthenticationProvider(private val accessTokenContext: AccessTokenContext) : AuthenticationProvider {
    override fun authenticate(authentication: Authentication): Authentication {
        authentication.isAuthenticated = true  //accessTokenContext.accessToken != null
        return authentication
    }

    override fun supports(authentication: Class<*>?) = true
}
