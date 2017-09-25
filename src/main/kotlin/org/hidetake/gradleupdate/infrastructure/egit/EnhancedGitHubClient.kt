package org.hidetake.gradleupdate.infrastructure.egit

import org.eclipse.egit.github.core.client.GitHubClient
import org.eclipse.egit.github.core.client.GitHubRequest
import org.eclipse.egit.github.core.client.GitHubResponse
import org.slf4j.LoggerFactory
import org.springframework.http.HttpHeaders
import java.net.HttpURLConnection

open class EnhancedGitHubClient(private val responseCacheRepository: ResponseCacheRepository) : GitHubClient() {
    private val log = LoggerFactory.getLogger(javaClass)

    fun setAccessToken(accessToken: String) {
        setOAuth2Token(accessToken)
    }

    override fun createConnection(uri: String, method: String): HttpURLConnection {
        log.debug("$method $uri")
        return super.createConnection(uri, method)
    }

    override fun get(request: GitHubRequest): GitHubResponse {
        val uri = request.generateUri()
        val httpRequest = createGet(uri)
        request.responseContentType?.also { accept ->
            httpRequest.setRequestProperty(HttpHeaders.ACCEPT, accept)
        }

        val responseCache = responseCacheRepository.find(uri)
        responseCache?.also {
            httpRequest.setRequestProperty(HttpHeaders.IF_NONE_MATCH, responseCache.eTag)
        }

        val code = httpRequest.responseCode
        updateRateLimits(httpRequest)
        return when {
            isOk(code) ->
                GitHubResponse(httpRequest, getBody(request, getStream(httpRequest))).also { response ->
                    val eTag = response.getHeader(HttpHeaders.ETAG)
                    responseCacheRepository.save(uri, ResponseCache(eTag, response.body))
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
