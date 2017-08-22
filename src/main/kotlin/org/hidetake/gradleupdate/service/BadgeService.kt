package org.hidetake.gradleupdate.service

import org.hidetake.gradleupdate.domain.GradleWrapper
import org.hidetake.gradleupdate.domain.GradleWrapperFactory
import org.kohsuke.github.GitHub
import org.springframework.stereotype.Service

@Service
class BadgeService(val github: GitHub) {
    fun findGradleWrapper(repositoryName: String): GradleWrapper? {
        val repository = github.getRepository(repositoryName)
        return GradleWrapperFactory.parseRepository(repository)
    }
}
