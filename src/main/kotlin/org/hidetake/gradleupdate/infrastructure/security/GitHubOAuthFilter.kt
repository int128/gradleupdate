package org.hidetake.gradleupdate.infrastructure.security

import org.springframework.security.web.authentication.preauth.AbstractPreAuthenticatedProcessingFilter
import javax.servlet.http.HttpServletRequest

class GitHubOAuthFilter : AbstractPreAuthenticatedProcessingFilter() {
    override fun getPreAuthenticatedPrincipal(request: HttpServletRequest) =
        request.cookies?.find { it.name == "S" }?.value

    override fun getPreAuthenticatedCredentials(request: HttpServletRequest) = ""
}
