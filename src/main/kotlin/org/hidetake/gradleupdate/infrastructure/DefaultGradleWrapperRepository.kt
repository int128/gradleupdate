package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.client.GitHubClient
import org.eclipse.egit.github.core.service.ContentsService
import org.hidetake.gradleupdate.domain.GradleWrapperRepository
import org.hidetake.gradleupdate.domain.GradleWrapperVersion
import org.springframework.stereotype.Repository
import java.util.*

@Repository
class DefaultGradleWrapperRepository(val client: GitHubClient) : GradleWrapperRepository {
    override fun findVersion(repositoryName: String): GradleWrapperVersion? =
        nullIfNotFound {
            ContentsService(client).getContents({repositoryName}, "gradle/wrapper/gradle-wrapper.properties")
        }
        ?.firstOrNull()
        ?.let { content ->
            val decoded = String(Base64.getMimeDecoder().decode(content.content))
            Regex("""distributionUrl=.+?/gradle-(.+?)-(.+?)\.zip""")
                .find(decoded)
                ?.groupValues
                ?.let { m -> GradleWrapperVersion(m[1], m[2]) }
        }
}
