package org.hidetake.gradleupdate.infrastructure

import org.hidetake.gradleupdate.infrastructure.egit.EnhancedGitHubClient
import org.hidetake.gradleupdate.infrastructure.egit.ResponseCacheRepository
import org.springframework.beans.factory.InitializingBean
import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.stereotype.Component

@Component
@ConfigurationProperties("gradleupdate.github")
class SystemGitHubClient(responseCacheRepository: ResponseCacheRepository)
    : EnhancedGitHubClient(responseCacheRepository), InitializingBean {
    private var accessToken: String? = null

    init {
        // EGit always sends author and committer on Commit API
        // but GitHub rejects null value.
        setSerializeNulls(false)
    }

    override fun getAccessToken() = accessToken!!

    fun setAccessToken(value: String) {
        accessToken = value
    }

    override fun afterPropertiesSet() {
        assert(accessToken != null)
    }
}
