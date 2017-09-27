package org.hidetake.gradleupdate.app

import com.google.appengine.api.memcache.Expiration
import com.google.appengine.api.memcache.MemcacheService
import org.hidetake.gradleupdate.infrastructure.egit.GitHubOAuthClient
import org.hidetake.gradleupdate.infrastructure.oauth.AccessTokenContext
import org.springframework.stereotype.Component
import java.util.*

@Component
class LoginService(
    private val client: GitHubOAuthClient,
    private val memcacheService: MemcacheService,
    private val context: AccessTokenContext
) {
    fun getRedirectURL() = client.redirectUrl

    fun createAuthorizationParameters(backTo: String) =
        client.computeAuthorizationParameters(backTo,
            UUID.randomUUID().toString().also { state ->
                memcacheService.put(state, true, Expiration.byDeltaSeconds(60))
            })

    fun continueAuthorization(state: String, code: String) {
        when (memcacheService.delete(state)) {
            true -> {
                val accessToken = client.acquireAccessToken(code)
                context.set(accessToken)
            }
            false ->
                throw IllegalStateException("OAuth state did not match: state=$state")
        }
    }

    fun logout() {
        context.clear()
    }
}
