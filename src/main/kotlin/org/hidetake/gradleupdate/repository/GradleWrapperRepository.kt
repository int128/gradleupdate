package org.hidetake.gradleupdate.repository

import org.hidetake.gradleupdate.domain.GradleWrapper
import org.hidetake.gradleupdate.domain.GradleWrapperFactory
import org.kohsuke.github.GitHub
import org.springframework.stereotype.Repository

@Repository
class GradleWrapperRepository(val gitHub: GitHub) {
    fun find(repositoryName: String): GradleWrapper? {
        val content = gitHub.getRepository(repositoryName)
            .getFileContent("gradle/wrapper/gradle-wrapper.properties")
            .read().bufferedReader().readText()
        return GradleWrapperFactory.parsePropertiesFile(content)
    }

    fun findLatestTemplate() = find("int128/latest-gradle-wrapper")
}
