package org.hidetake.gradleupdate.infrastructure.oauth

import java.io.Serializable

data class AccessToken(val value: String) : Serializable {
    companion object {
        const val serialVersionUID: Long = 1
    }
}
