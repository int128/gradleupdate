package org.hidetake.gradleupdate.infrastructure.oauth

import com.google.appengine.api.memcache.MemcacheService
import org.springframework.stereotype.Component
import org.springframework.web.filter.OncePerRequestFilter
import org.springframework.web.util.CookieGenerator
import java.util.*
import javax.servlet.FilterChain
import javax.servlet.http.HttpServletRequest
import javax.servlet.http.HttpServletResponse

@Component
class OAuthContextFilter(
    private val context: OAuthContext,
    private val memcacheService: MemcacheService
) : OncePerRequestFilter() {
    private val COOKIE = CookieGenerator().apply {
        cookieName = "S"
        cookiePath = "/"
        cookieMaxAge = 60 * 60 * 24 * 90
        isCookieHttpOnly = true
        isCookieSecure = true
    }

    private fun extractSessionId(request: HttpServletRequest) =
        request.cookies?.find { it.name == COOKIE.cookieName }?.let {
            try {
                UUID.fromString(it.value)
            } catch (e: IllegalArgumentException) {
                null
            }
        }

    override fun doFilterInternal(
        request: HttpServletRequest,
        response: HttpServletResponse,
        filterChain: FilterChain
    ) {
        val sessionId = extractSessionId(request)
        when (sessionId) {
            null -> {
                filterChain.doFilter(request, response)

                when (context.accessToken) {
                    null -> {}
                    else ->
                        UUID.randomUUID().also { uuid ->
                            memcacheService.put(uuid, context.accessToken)
                            COOKIE.addCookie(response, uuid.toString())
                        }
                }
            }
            else -> {
                val restoredAccessToken = memcacheService.get(sessionId) as AccessToken
                context.accessToken = restoredAccessToken

                filterChain.doFilter(request, response)

                when (context.accessToken) {
                    restoredAccessToken -> {}
                    null -> {
                        memcacheService.delete(sessionId)
                        COOKIE.removeCookie(response)
                    }
                    else -> {
                        memcacheService.put(sessionId, context.accessToken)
                    }
                }
            }
        }
    }
}
