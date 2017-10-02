package org.hidetake.gradleupdate

import com.google.appengine.api.memcache.MemcacheService
import com.google.appengine.api.memcache.MemcacheServiceFactory
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.web.servlet.ServletContextInitializer
import org.springframework.boot.web.support.SpringBootServletInitializer
import org.springframework.context.annotation.Bean

@SpringBootApplication
open class App : SpringBootServletInitializer() {
    @Bean
    open fun memcacheService(): MemcacheService =
        MemcacheServiceFactory.getMemcacheService()

    @Bean
    open fun sessionCookieConfigInitializer(): ServletContextInitializer =
        ServletContextInitializer { servletContext ->
            servletContext.sessionCookieConfig.apply {
                name = "S"
            }
        }
}
