package org.hidetake.gradleupdate

import org.eclipse.egit.github.core.client.GitHubClient
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.web.support.SpringBootServletInitializer
import org.springframework.context.annotation.Bean

@SpringBootApplication
open class App : SpringBootServletInitializer() {
    @Bean
    open fun gitHubClient(): GitHubClient = GitHubClient().apply {
    }
}