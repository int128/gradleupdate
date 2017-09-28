package org.hidetake.gradleupdate.infrastructure.security

import org.hidetake.gradleupdate.infrastructure.egit.AccessToken
import org.springframework.stereotype.Component
import org.springframework.web.context.annotation.SessionScope

@Component
@SessionScope
open class AccessTokenContext(var accessToken: AccessToken? = null)
