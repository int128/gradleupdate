package org.hidetake.gradleupdate.infrastructure.egit

import java.io.Serializable

class ResponseCache(val eTag: String, val body: Any) : Serializable {
    companion object {
        const val serialVersionUID: Long = 1
    }
}
