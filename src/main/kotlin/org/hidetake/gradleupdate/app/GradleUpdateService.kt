package org.hidetake.gradleupdate.app

import org.hidetake.gradleupdate.domain.GradleWrapperVersionStatus
import org.springframework.stereotype.Service

@Service
class GradleUpdateService(private val gitHub: GitHub) {
    fun getRepository(repositoryName: String) =
        gitHub.getRepository(repositoryName)

    fun getLatestGradleWrapperRepository() =
        getRepository("int128/latest-gradle-wrapper")

    fun getGradleWrapperVersionStatus(repositoryName: String): GradleWrapperVersionStatus? {
        val targetRepository = getRepository(repositoryName)
        val latestRepository = getLatestGradleWrapperRepository()
        return targetRepository.findGradleWrapperVersion()?.let { target ->
            latestRepository.findGradleWrapperVersion()?.let { latest ->
                GradleWrapperVersionStatus(target, latest)
            }
        }
    }
}
