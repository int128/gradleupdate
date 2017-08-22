package org.hidetake.gradleupdate

import org.kohsuke.github.GitHub
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.web.support.SpringBootServletInitializer
import org.springframework.context.annotation.Bean

@SpringBootApplication
open class App : SpringBootServletInitializer() {
    @Bean
    open fun gitHub(): GitHub = GitHub.connectAnonymously()
}
