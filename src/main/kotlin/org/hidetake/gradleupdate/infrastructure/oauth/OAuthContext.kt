package org.hidetake.gradleupdate.infrastructure.oauth

import org.springframework.stereotype.Component
import org.springframework.web.context.annotation.RequestScope

@Component
@RequestScope
open class OAuthContext(var accessToken: AccessToken? = null)
