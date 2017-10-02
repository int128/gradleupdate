package org.hidetake.gradleupdate.infrastructure

import org.hidetake.gradleupdate.infrastructure.egit.EnhancedGitHubClient
import org.hidetake.gradleupdate.infrastructure.egit.ResponseCacheRepository
import org.springframework.stereotype.Component

@Component
class SystemGitHubClient(responseCacheRepository: ResponseCacheRepository)
    : EnhancedGitHubClient(responseCacheRepository) {
    private val accessToken = System.getenv("SYSTEM_GITHUB_ACCESS_TOKEN")!!

    init {
        // EGit always sends author and committer on Commit API
        // but GitHub rejects null value.
        setSerializeNulls(false)
    }

    override fun getAccessToken() = accessToken
}
