package org.hidetake.gradleupdate

import com.google.appengine.api.memcache.MemcacheService
import com.google.appengine.api.memcache.MemcacheServiceFactory
import org.eclipse.egit.github.core.client.GitHubClient
import org.hidetake.gradleupdate.infrastructure.egit.EnhancedGitHubClient
import org.hidetake.gradleupdate.infrastructure.egit.ResponseCacheRepository
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.web.support.SpringBootServletInitializer
import org.springframework.context.annotation.Bean

@SpringBootApplication
open class App : SpringBootServletInitializer() {
    @Bean
    open fun gitHubClient(responseCacheRepository: ResponseCacheRepository): GitHubClient =
        EnhancedGitHubClient(responseCacheRepository).apply {
            setOAuth2Token(System.getenv("GITHUB_TOKEN"))
        }

    @Bean
    open fun memcacheService(): MemcacheService =
        MemcacheServiceFactory.getMemcacheService()
}
