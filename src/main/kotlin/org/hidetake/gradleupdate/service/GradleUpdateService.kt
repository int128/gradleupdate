package org.hidetake.gradleupdate.service

import org.hidetake.gradleupdate.domain.GradleWrapperStatus
import org.hidetake.gradleupdate.repository.GradleWrapperRepository
import org.springframework.stereotype.Service

@Service
class GradleUpdateService(val gradleWrapperRepository: GradleWrapperRepository) {
    fun getStatus(repositoryName: String): GradleWrapperStatus? {
        val current = gradleWrapperRepository.find(repositoryName)
        val latest = gradleWrapperRepository.findLatestTemplate()
        return current?.compareToLatest(latest)
    }
}
