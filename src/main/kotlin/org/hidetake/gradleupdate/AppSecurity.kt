package org.hidetake.gradleupdate

import org.springframework.boot.autoconfigure.security.oauth2.client.EnableOAuth2Sso
import org.springframework.context.annotation.Configuration
import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter

@EnableOAuth2Sso
@Configuration
open class AppSecurity : WebSecurityConfigurerAdapter() {
    override fun configure(http: HttpSecurity) {
        http.authorizeRequests()
            .antMatchers("/my", "/my/**").authenticated()
            .anyRequest().permitAll()

        http.logout()
            .logoutSuccessUrl("/")
    }
}
