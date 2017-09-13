package org.hidetake.gradleupdate.app

import com.google.appengine.api.memcache.Expiration
import com.google.appengine.api.memcache.MemcacheService
import org.eclipse.egit.github.core.client.GitHubClient
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.stereotype.Component
import java.util.*

@Component
@ConfigurationProperties("github.oauth")
class GitHubOAuthService(private val memcacheService: MemcacheService) {
    private val gitHubClient = GitHubClient()

    val authorizationEndpoint = "https://github.com/login/oauth/authorize"
    var clientId: String = ""
    var scope: String = ""

    fun buildAuthorizationParameters(backTo: String) = mapOf(
        "client_id" to clientId,
        "scope" to scope,
        "redirect_uri" to backTo,
        "state" to UUID.randomUUID().toString().also { state ->
            memcacheService.put(state, true, Expiration.byDeltaSeconds(5 * 60))
        }
    )

    fun exchangeCodeAndToken(state: String, code: String) =
        when (memcacheService.delete(state)) {
            true -> {
                // TODO
            }
            false ->
                throw IllegalStateException("OAuth state did not match: state=$state")
        }
}
