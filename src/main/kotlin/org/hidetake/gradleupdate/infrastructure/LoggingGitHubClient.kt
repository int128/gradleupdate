package org.hidetake.gradleupdate.infrastructure

import org.eclipse.egit.github.core.client.GitHubClient
import org.slf4j.LoggerFactory
import java.net.HttpURLConnection

val log = LoggerFactory.getLogger(LoggingGitHubClient::class.java)

class LoggingGitHubClient : GitHubClient() {
    override fun createConnection(uri: String, method: String): HttpURLConnection {
        log.debug("$method $uri")
        return super.createConnection(uri, method)
    }
}
