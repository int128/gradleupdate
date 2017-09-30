package org.hidetake.gradleupdate.infrastructure.egit

import org.eclipse.egit.github.core.client.GitHubClient
import org.eclipse.egit.github.core.client.GitHubRequest
import org.eclipse.egit.github.core.client.GitHubResponse
import org.slf4j.LoggerFactory
import org.springframework.http.HttpHeaders
import java.net.HttpURLConnection

class EnhancedGitHubClient(
    private val responseCacheRepository: ResponseCacheRepository,
    private val accessTokenProvider: () -> String
) : GitHubClient() {
    private val log = LoggerFactory.getLogger(javaClass)

    init {
        // EGit always sends author and committer on Commit API
        // but GitHub rejects null value.
        setSerializeNulls(false)
    }

    override fun setOAuth2Token(token: String?): GitHubClient =
        throw UnsupportedOperationException()

    override fun createConnection(uri: String, method: String): HttpURLConnection =
        super.createConnection(uri, method).also { connection ->
            log.debug("$method $uri")
            val accessToken = accessTokenProvider()
            connection.setRequestProperty("Authorization", "token $accessToken")
        }

    override fun get(request: GitHubRequest): GitHubResponse {
        val uri = request.generateUri()
        val httpRequest = createGet(uri)
        request.responseContentType?.also { accept ->
            httpRequest.setRequestProperty(HttpHeaders.ACCEPT, accept)
        }

        val requestProperties = httpRequest.requestProperties
        val responseCache = responseCacheRepository.find(ResponseCacheKey(uri, requestProperties))
        responseCache?.also {
            httpRequest.setRequestProperty(HttpHeaders.IF_NONE_MATCH, responseCache.eTag)
        }

        val code = httpRequest.responseCode
        updateRateLimits(httpRequest)
        return when {
            isOk(code) ->
                GitHubResponse(httpRequest, getBody(request, getStream(httpRequest))).also { response ->
                    val eTag = response.getHeader(HttpHeaders.ETAG)
                    responseCacheRepository.save(
                        ResponseCacheKey(uri, requestProperties),
                        ResponseCache(eTag, response.body))
                    log.debug("CACHED {} @ {}", uri, eTag)
                }
            isEmpty(code) ->
                GitHubResponse(httpRequest, null)
            isNotModified(code) ->
                GitHubResponse(httpRequest, responseCache?.body).also { response ->
                    val eTag = response.getHeader(HttpHeaders.ETAG)
                    log.debug("HIT {} @ {}", uri, eTag)
                }
            else ->
                throw createException(getStream(httpRequest), code, httpRequest.responseMessage)
        }
    }

    private fun isNotModified(code: Int) = code == 304
}
