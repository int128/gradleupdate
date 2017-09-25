package org.hidetake.gradleupdate

import com.google.appengine.api.memcache.MemcacheService
import com.google.appengine.api.memcache.MemcacheServiceFactory
import org.hidetake.gradleupdate.infrastructure.egit.EnhancedGitHubClient
import org.hidetake.gradleupdate.infrastructure.egit.GitHubOAuthClient
import org.hidetake.gradleupdate.infrastructure.egit.ResponseCacheRepository
import org.hidetake.gradleupdate.infrastructure.oauth.OAuthContext
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.web.support.SpringBootServletInitializer
import org.springframework.context.annotation.Bean
import org.springframework.web.context.annotation.RequestScope

@SpringBootApplication
open class App : SpringBootServletInitializer() {
    @Bean
    @ConfigurationProperties("gradleupdate.github")
    open fun gradleUpdateGitHubClient(responseCacheRepository: ResponseCacheRepository) =
        EnhancedGitHubClient(responseCacheRepository).apply {
            // EGit always sends author and committer on Commit API
            // but GitHub rejects null value.
            setSerializeNulls(false)
        }

    @Bean
    @RequestScope
    open fun contextGitHubClient(responseCacheRepository: ResponseCacheRepository, context: OAuthContext) =
        EnhancedGitHubClient(responseCacheRepository).apply {
            context.accessToken?.also { setAccessToken(it.value) }
        }

    @Bean
    @ConfigurationProperties("github.oauth")
    open fun gitHubOAuthClient() = GitHubOAuthClient()

    @Bean
    open fun memcacheService(): MemcacheService =
        MemcacheServiceFactory.getMemcacheService()
}
