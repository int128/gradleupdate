package org.hidetake.gradleupdate.app

import com.google.appengine.api.memcache.Expiration
import com.google.appengine.api.memcache.MemcacheService
import org.hidetake.gradleupdate.infrastructure.egit.AccessToken
import org.hidetake.gradleupdate.infrastructure.egit.GitHubOAuthClient
import org.springframework.stereotype.Component
import java.util.*

@Component
class LoginService(
    private val client: GitHubOAuthClient,
    private val memcacheService: MemcacheService
) {
    fun getRedirectURL() = client.redirectUrl

    fun createAuthorizationParameters(backTo: String) =
        client.computeAuthorizationParameters(backTo,
            UUID.randomUUID().toString().also { state ->
                memcacheService.put(state, true, Expiration.byDeltaSeconds(60))
            })

    fun continueAuthorization(state: String, code: String): AccessToken =
        when (memcacheService.delete(state)) {
            true ->
                client.acquireAccessToken(code)
            false ->
                throw IllegalStateException("OAuth state did not match: state=$state")
        }
}
