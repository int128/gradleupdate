package org.hidetake.gradleupdate

import com.google.appengine.api.memcache.MemcacheService
import com.google.appengine.api.memcache.MemcacheServiceFactory
import org.hidetake.gradleupdate.infrastructure.egit.EnhancedGitHubClient
import org.hidetake.gradleupdate.infrastructure.egit.ResponseCacheRepository
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.web.support.SpringBootServletInitializer
import org.springframework.context.annotation.Bean
import org.springframework.security.oauth2.client.OAuth2ClientContext

@SpringBootApplication
open class App : SpringBootServletInitializer() {
    @Bean
    open fun gradleUpdateGitHubClient(responseCacheRepository: ResponseCacheRepository): EnhancedGitHubClient {
        val accessToken = System.getenv("GRADLEUPDATE_GITHUB_ACCESS_TOKEN")
        return EnhancedGitHubClient(responseCacheRepository, { accessToken })
    }

    @Bean
    open fun contextGitHubClient(responseCacheRepository: ResponseCacheRepository, clientContext: OAuth2ClientContext): EnhancedGitHubClient =
        EnhancedGitHubClient(responseCacheRepository, { clientContext.accessToken.value })

    @Bean
    open fun memcacheService(): MemcacheService =
        MemcacheServiceFactory.getMemcacheService()
}
