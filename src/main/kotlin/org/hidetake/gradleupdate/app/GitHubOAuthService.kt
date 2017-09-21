package org.hidetake.gradleupdate.app

import com.google.appengine.api.memcache.Expiration
import com.google.appengine.api.memcache.MemcacheService
import org.hidetake.gradleupdate.infrastructure.egit.GitHubOAuthClient
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.stereotype.Component
import java.util.*

@Component
@ConfigurationProperties("github.oauth")
class GitHubOAuthService(private val memcacheService: MemcacheService) {
    val authorizationEndpoint = "https://github.com/login/oauth/authorize"

    var clientId: String = ""
    var clientSecret: String = ""
    var scope: String = ""

    fun createAuthorizationParameters(backTo: String) = mapOf(
        "client_id" to clientId,
        "scope" to scope,
        "redirect_uri" to backTo,
        "state" to UUID.randomUUID().toString().also { state ->
            memcacheService.put(state, true, Expiration.byDeltaSeconds(5 * 60))
        }
    )

    fun createSession(state: String, code: String): String =
        when (memcacheService.delete(state)) {
            true -> {
                val client = GitHubOAuthClient(clientId, clientSecret)
                val accessToken = client.acquireAccessToken(code)
                UUID.randomUUID().toString().also { sessionId ->
                    memcacheService.put(sessionId, accessToken)
                }
            }
            false ->
                throw IllegalStateException("OAuth state did not match: state=$state")
        }
}
