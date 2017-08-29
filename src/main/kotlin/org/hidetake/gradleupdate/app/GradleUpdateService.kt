package org.hidetake.gradleupdate.app

import org.hidetake.gradleupdate.domain.GradleWrapperRepository
import org.hidetake.gradleupdate.domain.GradleWrapperVersionStatus
import org.hidetake.gradleupdate.domain.RepositoryRepository
import org.springframework.stereotype.Service

@Service
class GradleUpdateService(
    private val repositoryRepository: RepositoryRepository,
    private val gradleWrapperRepository: GradleWrapperRepository
) {
    fun getRepository(repositoryName: String) =
        repositoryRepository.getByName(repositoryName)

    fun getGradleWrapperVersionStatus(repositoryName: String): GradleWrapperVersionStatus? =
        gradleWrapperRepository.findVersion(repositoryName)?.let { target ->
            gradleWrapperRepository.findVersion("int128/latest-gradle-wrapper")?.let { latest ->
                GradleWrapperVersionStatus(target, latest)
            }
        }
}
