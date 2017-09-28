package org.hidetake.gradleupdate

import com.google.appengine.api.memcache.MemcacheService
import com.google.appengine.api.memcache.MemcacheServiceFactory
import org.hidetake.gradleupdate.infrastructure.egit.AccessToken
import org.hidetake.gradleupdate.infrastructure.egit.EnhancedGitHubClient
import org.hidetake.gradleupdate.infrastructure.egit.GitHubOAuthClient
import org.hidetake.gradleupdate.infrastructure.egit.ResponseCacheRepository
import org.hidetake.gradleupdate.infrastructure.session.MemcacheSession
import org.hidetake.gradleupdate.infrastructure.session.MemcacheSessionRepository
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.web.support.SpringBootServletInitializer
import org.springframework.context.annotation.Bean
import org.springframework.session.web.http.SessionRepositoryFilter
import javax.servlet.ServletContext


@SpringBootApplication
open class App : SpringBootServletInitializer() {
    @Bean
    open fun gradleUpdateGitHubClient(responseCacheRepository: ResponseCacheRepository): EnhancedGitHubClient {
        val accessToken = AccessToken(System.getenv("GRADLEUPDATE_GITHUB_ACCESS_TOKEN"))
        return EnhancedGitHubClient(responseCacheRepository, { accessToken })
    }

    @Bean
    open fun contextGitHubClient(responseCacheRepository: ResponseCacheRepository): EnhancedGitHubClient =
        EnhancedGitHubClient(responseCacheRepository, {
            throw IllegalStateException("Login required")
        })

    @Bean
    open fun gitHubOAuthClient(): GitHubOAuthClient =
        GitHubOAuthClient(
            System.getenv("GITHUB_OAUTH_CLIENT_ID"),
            System.getenv("GITHUB_OAUTH_CLIENT_SECRET"),
            "public_repo",
            "https://github.com/login/oauth/authorize"
        )

    @Bean
    open fun memcacheService(): MemcacheService =
        MemcacheServiceFactory.getMemcacheService()

    @Bean
    open fun sessionRepositoryFilter(repository: MemcacheSessionRepository, servletContext: ServletContext) =
        SessionRepositoryFilter<MemcacheSession>(repository).apply {
            setServletContext(servletContext)
        }
}
