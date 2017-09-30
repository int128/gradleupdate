package org.hidetake.gradleupdate.infrastructure.egit

import java.io.Serializable

data class ResponseCacheKey(
    val uri: String,
    val requestProperties: Map<String, List<String>>
) : Serializable {
    companion object {
        const val serialVersionUID: Long = 1
    }
}
