package org.hidetake.gradleupdate.infrastructure.security

import org.springframework.boot.autoconfigure.security.oauth2.resource.UserInfoRestTemplateCustomizer
import org.springframework.security.oauth2.client.OAuth2RestTemplate
import org.springframework.stereotype.Component

@Component
class GitHubUserInfoRestTemplateCustomizer(
    private val requestInterceptor: GitHubUserInfoRestRequestInterceptor
) : UserInfoRestTemplateCustomizer {
    override fun customize(template: OAuth2RestTemplate) {
        template.interceptors.add(requestInterceptor)
    }
}
