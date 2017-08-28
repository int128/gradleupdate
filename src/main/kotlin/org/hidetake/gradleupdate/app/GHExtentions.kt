package org.hidetake.gradleupdate.app

import org.hidetake.gradleupdate.domain.GradleWrapperVersion

fun GHRepository.findGradleWrapperVersion(): GradleWrapperVersion? {
    val content = getFileContent("gradle/wrapper/gradle-wrapper.properties")
        .read()
        .bufferedReader()
        .readText()
    return Regex("""distributionUrl=.+?/gradle-(.+?)-(.+?)\.zip""").find(content)
        ?.groupValues
        ?.let { m -> GradleWrapperVersion(m[1], m[2]) }
}
