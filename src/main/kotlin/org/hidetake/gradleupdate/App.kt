package org.hidetake.gradleupdate

import com.google.appengine.api.memcache.MemcacheService
import com.google.appengine.api.memcache.MemcacheServiceFactory
import org.springframework.boot.SpringApplication
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.web.servlet.support.SpringBootServletInitializer
import org.springframework.context.annotation.Bean

@SpringBootApplication
open class App : SpringBootServletInitializer() {
    @Bean
    open fun memcacheService(): MemcacheService = MemcacheServiceFactory.getMemcacheService()
}

fun main(args: Array<String>) {
    SpringApplication.run(App::class.java, *args)
}
