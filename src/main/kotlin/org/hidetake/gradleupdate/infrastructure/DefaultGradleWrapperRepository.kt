package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.client.GitHubClient
import org.eclipse.egit.github.core.service.ContentsService
import org.hidetake.gradleupdate.domain.GradleWrapperFile
import org.hidetake.gradleupdate.domain.GradleWrapperRepository
import org.hidetake.gradleupdate.domain.GradleWrapperVersion
import org.springframework.stereotype.Repository
import java.util.*

@Repository
class DefaultGradleWrapperRepository(client: GitHubClient) : GradleWrapperRepository {
    private val contentsService = ContentsService(client)

    override fun findVersion(repositoryName: String): GradleWrapperVersion? =
        findFile(repositoryName, "gradle/wrapper/gradle-wrapper.properties")
        ?.let { content ->
            val decoded = String(Base64.getMimeDecoder().decode(content.content))
            Regex("""distributionUrl=.+?/gradle-(.+?)-(.+?)\.zip""")
                .find(decoded)
                ?.groupValues
                ?.let { m -> GradleWrapperVersion(m[1], m[2]) }
        }

    override fun findFiles(repositoryName: String): List<GradleWrapperFile> =
        listOf(
            GradleWrapperFile("gradle/wrapper/gradle-wrapper.properties"),
            GradleWrapperFile("gradle/wrapper/gradle-wrapper.jar"),
            GradleWrapperFile("gradlew", true),
            GradleWrapperFile("gradlew.bat")
        ).map { file ->
            GradleWrapperFile(file.path, file.executable, findFile(repositoryName, file.path)?.content)
        }

    private fun findFile(repositoryName: String, path: String) =
        nullIfNotFound {
            contentsService.getContents({repositoryName}, path)
        }?.firstOrNull()
}
