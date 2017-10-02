package org.hidetake.gradleupdate.infrastructure

import org.hidetake.gradleupdate.infrastructure.egit.EnhancedGitHubClient
import org.hidetake.gradleupdate.infrastructure.egit.ResponseCacheRepository
import org.springframework.security.oauth2.client.OAuth2ClientContext
import org.springframework.stereotype.Component

@Component
class LoginUserGitHubClient(
    private val context: OAuth2ClientContext,
    responseCacheRepository: ResponseCacheRepository
) : EnhancedGitHubClient(responseCacheRepository) {
    override fun getAccessToken(): String =
        context.accessToken?.value ?: throw IllegalStateException("Login required")
}
