package org.hidetake.gradleupdate.domain

import org.kohsuke.github.GHRepository

object GradleWrapperFactory {
    fun parseRepository(repository: GHRepository): GradleWrapper? {
        val gradleWrapperProperties = readContent(repository, "gradle/wrapper/gradle-wrapper.properties")
        return gradleWrapperProperties?.let { parsePropertiesFile(it) }
    }

    fun parsePropertiesFile(content: String): GradleWrapper? {
        val match = Regex("""distributionUrl=.+?/gradle-(.+?)-.+?\.zip""").find(content)
        return match?.groups?.get(1)?.value?.let { GradleWrapper(it) }
    }

    private fun readContent(repository: GHRepository, path: String): String? {
        return repository.getFileContent(path).read().bufferedReader().readText()
    }
}
