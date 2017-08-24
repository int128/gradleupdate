package org.hidetake.gradleupdate.service

import org.hidetake.gradleupdate.domain.GradleWrapperVersionStatus
import org.kohsuke.github.GHRepository
import org.kohsuke.github.GitHub
import org.springframework.stereotype.Service

@Service
class GradleUpdateService(private val gitHub: GitHub) {
    fun getGradleWrapperVersionStatus(repositoryName: String) =
        getGradleWrapperVersionStatus(gitHub.getRepository(repositoryName))

    fun getGradleWrapperVersionStatus(targetRepository: GHRepository) =
        getGradleWrapperVersionStatus(targetRepository, gitHub.getRepository("int128/latest-gradle-wrapper"))

    fun getGradleWrapperVersionStatus(targetRepository: GHRepository, latestRepository: GHRepository) =
        targetRepository.findGradleWrapperVersion()?.let { target ->
            latestRepository.findGradleWrapperVersion()?.let { latest ->
                GradleWrapperVersionStatus(target, latest)
            }
        }

    fun getRepository(repositoryName: String) = gitHub.getRepository(repositoryName)
}
