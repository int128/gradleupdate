package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.service.ContentsService
import org.hidetake.gradleupdate.domain.GradleWrapperFile
import org.hidetake.gradleupdate.domain.GradleWrapperRepository
import org.hidetake.gradleupdate.domain.GradleWrapperVersion
import org.hidetake.gradleupdate.domain.RepositoryPath
import org.springframework.stereotype.Repository
import java.util.*

@Repository
class DefaultGradleWrapperRepository(client: SystemGitHubClient) : GradleWrapperRepository {
    private val contentsService = ContentsService(client)

    override fun findVersion(repositoryPath: RepositoryPath): GradleWrapperVersion? =
        findFile(repositoryPath, "gradle/wrapper/gradle-wrapper.properties")
        ?.let { content ->
            val decoded = String(Base64.getMimeDecoder().decode(content.content))
            Regex("""distributionUrl=.+?/gradle-(.+?)-.+?\.zip""")
                .find(decoded)
                ?.groupValues
                ?.let { m -> GradleWrapperVersion(m[1]) }
        }

    override fun findFiles(repositoryPath: RepositoryPath): List<GradleWrapperFile> =
        listOf(
            GradleWrapperFile("gradle/wrapper/gradle-wrapper.properties"),
            GradleWrapperFile("gradle/wrapper/gradle-wrapper.jar"),
            GradleWrapperFile("gradlew", true),
            GradleWrapperFile("gradlew.bat")
        ).map { file ->
            GradleWrapperFile(file.path, file.executable, findFile(repositoryPath, file.path)?.content)
        }

    private fun findFile(repositoryPath: RepositoryPath, path: String) =
        nullIfNotFound {
            contentsService.getContents({repositoryPath.fullName}, path)
        }?.firstOrNull()
}
