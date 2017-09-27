package org.hidetake.gradleupdate.infrastructure.oauth

import com.google.appengine.api.memcache.MemcacheService
import org.springframework.stereotype.Component
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import org.springframework.web.util.CookieGenerator
import java.util.*
import javax.servlet.http.HttpServletRequest

@Component
class AccessTokenContext(private val memcacheService: MemcacheService) {
    private val CACHE_ATTRIBUTE = javaClass.name

    private fun cookie(request: HttpServletRequest) = CookieGenerator().apply {
        cookieName = "S"
        cookiePath = "/"
        cookieMaxAge = 60 * 60 * 24 * 90
        isCookieHttpOnly = true
        isCookieSecure = request.isSecure
    }

    fun get(): AccessToken? {
        val context = getContext()
        return withCache(context.request) {
            extractSessionId(context.request)?.let { sessionId ->
                memcacheService.get(sessionId) as? AccessToken
            }
        }
    }

    fun set(accessToken: AccessToken) {
        val context = getContext()
        val sessionId = extractSessionId(context.request) ?: UUID.randomUUID().toString().also {
            cookie(context.request).addCookie(context.response, it)
        }
        memcacheService.put(sessionId, accessToken)
    }

    fun clear() {
        val context = getContext()
        extractSessionId(context.request)?.also { sessionId ->
            memcacheService.delete(sessionId)
            cookie(context.request).removeCookie(context.response)
        }
    }

    private class NoAccessToken

    private fun withCache(request: HttpServletRequest, fetch: () -> AccessToken?): AccessToken? {
        val cache = request.getAttribute(CACHE_ATTRIBUTE)
        return when (cache) {
            is AccessToken -> cache
            is NoAccessToken -> null
            else -> fetch().also { accessToken ->
                when (accessToken) {
                    is AccessToken -> request.setAttribute(CACHE_ATTRIBUTE, accessToken)
                    else -> request.setAttribute(CACHE_ATTRIBUTE, NoAccessToken())
                }
            }
        }
    }

    private fun extractSessionId(request: HttpServletRequest): String? =
        request.cookies?.find { it.name == "S" }?.value

    private fun getContext(): ServletRequestAttributes =
        (RequestContextHolder.getRequestAttributes() as? ServletRequestAttributes)
            ?: throw IllegalStateException("HttpServletRequest did not found")
}
