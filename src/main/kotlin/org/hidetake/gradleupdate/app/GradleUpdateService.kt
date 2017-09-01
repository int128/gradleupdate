package org.hidetake.gradleupdate.app

import org.hidetake.gradleupdate.domain.*
import org.springframework.stereotype.Service

@Service
class GradleUpdateService(
    private val repositoryRepository: RepositoryRepository,
    private val gradleWrapperRepository: GradleWrapperRepository,
    private val pullRequestRepository: PullRequestRepository
) {
    private val LATEST_GRADLE_WRAPPER = "int128/latest-gradle-wrapper"

    fun getRepository(repositoryName: String) =
        repositoryRepository.getByName(repositoryName)

    fun getGradleWrapperVersionStatus(repositoryName: String): GradleWrapperVersionStatus? =
        gradleWrapperRepository.findVersion(repositoryName)?.let { target ->
            gradleWrapperRepository.findVersion(LATEST_GRADLE_WRAPPER)?.let { latest ->
                GradleWrapperVersionStatus(target, latest)
            }
        }

    fun createPullRequestForLatestGradleWrapper(repositoryName: String) =
        gradleWrapperRepository.findVersion(repositoryName)?.let { target ->
            gradleWrapperRepository.findVersion(LATEST_GRADLE_WRAPPER)?.let { latest ->
                val status = GradleWrapperVersionStatus(target, latest)
                when {
                    status.upToDate -> TODO()
                    else -> {
                        val files = gradleWrapperRepository.findFiles(LATEST_GRADLE_WRAPPER)
                        val pullRequest = GradleWrapperPullRequestFactory.create(repositoryName, latest, files)
                        pullRequestRepository.createOrUpdate(pullRequest)
                    }
                }
            }
        }
}
