package org.hidetake.gradleupdate

import com.google.appengine.api.memcache.MemcacheService
import com.google.appengine.api.memcache.MemcacheServiceFactory
import org.hidetake.gradleupdate.infrastructure.session.MemcacheSessionRepository
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.web.support.SpringBootServletInitializer
import org.springframework.context.annotation.Bean
import org.springframework.session.web.http.SessionRepositoryFilter

@SpringBootApplication
open class App : SpringBootServletInitializer() {
    @Bean
    open fun memcacheService(): MemcacheService =
        MemcacheServiceFactory.getMemcacheService()

    @Bean
    open fun springSessionRepositoryFilter(memcacheSessionRepository: MemcacheSessionRepository) =
        SessionRepositoryFilter(memcacheSessionRepository)
}
