package org.hidetake.gradleupdate.infrastructure.egit

import org.eclipse.egit.github.core.client.GitHubClient
import org.springframework.beans.factory.InitializingBean

class GitHubOAuthClient : GitHubClient("github.com"), InitializingBean {
    var clientId: String? = null
    var clientSecret: String? = null
    var scope: String? = null
    var redirectUrl: String? = null

    init {
        headerAccept = "application/json"
    }

    override fun afterPropertiesSet() {
        assert(clientId != null)
        assert(clientSecret != null)
        assert(scope != null)
        assert(redirectUrl != null)
    }

    override fun configureUri(uri: String) = uri

    fun computeAuthorizationParameters(redirectUri: String, state: String) = mapOf(
        "client_id" to clientId,
        "scope" to scope,
        "redirect_uri" to redirectUri,
        "state" to state
    )

    fun acquireAccessToken(code: String): String =
        postAccessToken(code).let { (accessToken, error, errorDescription) ->
            accessToken ?: throw IllegalStateException(
                "Could not acquire access token: $error, $errorDescription")
        }

    private fun postAccessToken(code: String): AccessTokenResponse =
        post("/login/oauth/access_token", mapOf(
            "client_id" to clientId,
            "client_secret" to clientSecret,
            "code" to code
        ), AccessTokenResponse::class.java)

    private data class AccessTokenResponse(
        val accessToken: String?,
        val error: String?,
        val errorDescription: String?
    )
}
