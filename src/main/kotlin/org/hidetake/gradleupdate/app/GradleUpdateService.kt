package org.hidetake.gradleupdate.app

import org.eclipse.egit.github.core.PullRequest
import org.hidetake.gradleupdate.domain.*
import org.springframework.stereotype.Service

@Service
class GradleUpdateService(
    private val repositoryRepository: RepositoryRepository,
    private val gradleWrapperRepository: GradleWrapperRepository,
    private val pullRequestRepository: PullRequestRepository
) {
    private val LATEST_GRADLE_WRAPPER = RepositoryPath("int128", "latest-gradle-wrapper")

    fun getRepository(repositoryPath: RepositoryPath) =
        repositoryRepository.getByName(repositoryPath)

    fun getGradleWrapperVersionStatus(repositoryPath: RepositoryPath): GradleWrapperVersionStatus? =
        gradleWrapperRepository.findVersion(repositoryPath)?.let { target ->
            gradleWrapperRepository.findVersion(LATEST_GRADLE_WRAPPER)?.let { latest ->
                GradleWrapperVersionStatus(target, latest)
            }
        }

    fun getPullRequestStatus(repositoryPath: RepositoryPath): GradleWrapperPullRequestStatus = TODO()

    fun findPullRequestForLatestGradleWrapper(repositoryPath: RepositoryPath): PullRequest? =
        gradleWrapperRepository.findVersion(LATEST_GRADLE_WRAPPER)?.let { latest ->
            pullRequestRepository.find(repositoryPath, latest)
        }

    fun createPullRequestForLatestGradleWrapper(repositoryPath: RepositoryPath) =
        gradleWrapperRepository.findVersion(repositoryPath)?.let { target ->
            gradleWrapperRepository.findVersion(LATEST_GRADLE_WRAPPER)?.let { latest ->
                val status = GradleWrapperVersionStatus(target, latest)
                when {
                    status.upToDate -> TODO()
                    else -> {
                        val files = gradleWrapperRepository.findFiles(LATEST_GRADLE_WRAPPER)
                        pullRequestRepository.createOrUpdate(repositoryPath, latest, files)
                    }
                }
            }
        }
}
