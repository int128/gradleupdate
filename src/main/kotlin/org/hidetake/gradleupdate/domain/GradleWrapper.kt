package org.hidetake.gradleupdate.domain

class GradleWrapper(
    val version: String
) {
    fun compareToLatest(latest: GradleWrapper?): GradleWrapperStatus? =
        latest?.let { GradleWrapperStatus(this, version == latest.version) }
}
