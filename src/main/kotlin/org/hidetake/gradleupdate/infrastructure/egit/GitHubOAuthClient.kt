package org.hidetake.gradleupdate.infrastructure.egit

import org.eclipse.egit.github.core.client.GitHubClient

class GitHubOAuthClient(
    private val clientId: String,
    private val clientSecret: String
) : GitHubClient("github.com") {
    init {
        headerAccept = "application/json"
    }

    override fun configureUri(uri: String) = uri

    fun acquireAccessToken(code: String): String =
        postAccessToken(code).let { (accessToken, error, errorDescription) ->
            accessToken ?: throw IllegalStateException(
                "Could not acquire access token: $error, $errorDescription")
        }

    private data class AccessTokenResponse(
        val accessToken: String?,
        val error: String?,
        val errorDescription: String?
    )

    private fun postAccessToken(code: String): AccessTokenResponse =
        post("/login/oauth/access_token", mapOf(
            "client_id" to clientId,
            "client_secret" to clientSecret,
            "code" to code
        ), AccessTokenResponse::class.java)
}
