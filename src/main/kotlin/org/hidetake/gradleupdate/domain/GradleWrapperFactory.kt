package org.hidetake.gradleupdate.domain

object GradleWrapperFactory {
    fun parsePropertiesFile(content: String): GradleWrapper? =
        Regex("""distributionUrl=.+?/gradle-(.+?)-.+?\.zip""")
            .find(content)?.groups?.get(1)?.value?.let { GradleWrapper(it) }
}
