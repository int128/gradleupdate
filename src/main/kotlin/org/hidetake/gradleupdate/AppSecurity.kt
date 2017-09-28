package org.hidetake.gradleupdate

import com.google.appengine.api.memcache.MemcacheService
import org.hidetake.gradleupdate.infrastructure.security.GitHubOAuthFilter
import org.hidetake.gradleupdate.infrastructure.security.GitHubUserDetailsService
import org.springframework.context.annotation.Bean
import org.springframework.security.authentication.AccountStatusUserDetailsChecker
import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter
import org.springframework.security.web.authentication.preauth.PreAuthenticatedAuthenticationProvider

@EnableWebSecurity
open class AppSecurity : WebSecurityConfigurerAdapter() {
    override fun configure(http: HttpSecurity) {
        http.authorizeRequests()
            .antMatchers("/my", "/my/**").authenticated()
            .anyRequest().permitAll()

        http.addFilter(GitHubOAuthFilter().apply {
            setAuthenticationManager(authenticationManager())
        })
    }

    @Bean
    open fun authenticationProvider(memcacheService: MemcacheService) =
        PreAuthenticatedAuthenticationProvider().apply {
            setPreAuthenticatedUserDetailsService(GitHubUserDetailsService(memcacheService))
            setUserDetailsChecker(AccountStatusUserDetailsChecker())
        }
}
