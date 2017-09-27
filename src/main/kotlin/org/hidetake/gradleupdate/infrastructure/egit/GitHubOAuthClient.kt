package org.hidetake.gradleupdate.infrastructure.egit

import org.eclipse.egit.github.core.client.GitHubClient
import org.hidetake.gradleupdate.infrastructure.oauth.AccessToken

class GitHubOAuthClient(
    private val clientId: String,
    private val clientSecret: String,
    private val scope: String,
    val redirectUrl: String
) : GitHubClient("github.com") {
    init {
        headerAccept = "application/json"
    }

    override fun configureUri(uri: String) = uri

    fun computeAuthorizationParameters(redirectUri: String, state: String) = mapOf(
        "client_id" to clientId,
        "scope" to scope,
        "redirect_uri" to redirectUri,
        "state" to state
    )

    fun acquireAccessToken(code: String): AccessToken =
        postAccessToken(code).let { (accessToken, error, errorDescription) ->
            accessToken?.let {
                AccessToken(it)
            } ?: throw IllegalStateException(
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
